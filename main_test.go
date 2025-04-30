package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// Helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestLeerTransaccionesDesdeCSV_Success(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.csv")
	data := "ID,Tipo,Monto\n" +
		"1,Crédito,100.50\n" +
		"2,débito,50.25\n"
	if err := os.WriteFile(file, []byte(data), 0644); err != nil {
		t.Fatalf("No se pudo crear el archivo de prueba: %v", err)
	}

	trans, err := leerTransaccionesDesdeCSV(file)
	if err != nil {
		t.Fatalf("Se esperaba sin error, se obtuvo: %v", err)
	}
	if len(trans) != 2 {
		t.Errorf("Se esperaba 2 transacciones, se obtuvo %d", len(trans))
	}
	if trans[0].ID != "1" || trans[0].Tipo != "Crédito" || trans[0].Monto != 100.50 {
		t.Errorf("Transacción 1 mal parseada: %+v", trans[0])
	}
	if trans[1].ID != "2" || trans[1].Tipo != "débito" || trans[1].Monto != 50.25 {
		t.Errorf("Transacción 2 mal parseada: %+v", trans[1])
	}
}

func TestLeerTransaccionesDesdeCSV_FileNotFound(t *testing.T) {
	_, err := leerTransaccionesDesdeCSV("no_existe.csv")
	if err == nil {
		t.Error("Se esperaba un error al leer un archivo inexistente, pero err es nil")
	}
}

func TestLeerTransaccionesDesdeCSV_InvalidMonto(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "invalid.csv")
	data := "ID,Tipo,Monto\n" +
		"1,Crédito,abc\n"
	os.WriteFile(file, []byte(data), 0644)

	_, err := leerTransaccionesDesdeCSV(file)
	if err == nil {
		t.Error("Se esperaba un error al convertir monto inválido, pero err es nil")
	}
}

func TestLeerTransaccionesDesdeCSV_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "empty.csv")
	data := "ID,Tipo,Monto\n"
	if err := os.WriteFile(file, []byte(data), 0644); err != nil {
		t.Fatalf("No se pudo crear el archivo vacío: %v", err)
	}

	trans, err := leerTransaccionesDesdeCSV(file)
	if err != nil {
		t.Fatalf("Se esperaba sin error para archivo vacío, se obtuvo: %v", err)
	}
	if len(trans) != 0 {
		t.Errorf("Se esperaba 0 transacciones, se obtuvo %d", len(trans))
	}
}

func TestGenerarReporte(t *testing.T) {
	transacciones := []Transaccion{
		{ID: "1", Tipo: "Crédito", Monto: 100},
		{ID: "2", Tipo: "débito", Monto: 50},
		{ID: "3", Tipo: "credito", Monto: 200},
	}

	output := captureOutput(func() {
		generarReporte(transacciones)
	})

	if !bytes.Contains([]byte(output), []byte("Balance Final: 250.00")) {
		t.Errorf("Balance incorrecto en output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Transacción de Mayor Monto: ID 3 - 200.00")) {
		t.Errorf("Mayor monto incorrecto en output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Crédito: 2")) || !bytes.Contains([]byte(output), []byte("Débito: 1")) {
		t.Errorf("Conteo de transacciones incorrecto en output: %s", output)
	}
}

func TestGenerarReporte_Empty(t *testing.T) {
	output := captureOutput(func() {
		generarReporte(nil)
	})

	if !bytes.Contains([]byte(output), []byte("Balance Final: 0.00")) {
		t.Errorf("Balance incorrecto para slice vacío: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Transacción de Mayor Monto: ID  - 0.00")) {
		t.Errorf("Mayor monto incorrecto para slice vacío: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Crédito: 0 Débito: 0")) {
		t.Errorf("Conteo incorrecto para slice vacío: %s", output)
	}
}

func TestGenerarReporte_UnknownType(t *testing.T) {
	transacciones := []Transaccion{
		{ID: "1", Tipo: "Otro", Monto: 500},
		{ID: "2", Tipo: "Crédito", Monto: 100},
	}

	output := captureOutput(func() {
		generarReporte(transacciones)
	})

	if !bytes.Contains([]byte(output), []byte("Balance Final: 100.00")) {
		t.Errorf("Balance incorrecto con tipo desconocido: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Transacción de Mayor Monto: ID 1 - 500.00")) {
		t.Errorf("Mayor monto incorrecto con tipo desconocido: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Crédito: 1")) || !bytes.Contains([]byte(output), []byte("Débito: 0")) {
		t.Errorf("Conteo incorrecto con tipo desconocido: %s", output)
	}
}

func TestGenerarReporte_CaseInsensitive(t *testing.T) {
	transacciones := []Transaccion{
		{ID: "1", Tipo: "CRÉDITO", Monto: 150},
		{ID: "2", Tipo: "DEBITO", Monto: 50},
	}

	output := captureOutput(func() {
		generarReporte(transacciones)
	})

	if !bytes.Contains([]byte(output), []byte("Balance Final: 100.00")) {
		t.Errorf("Balance incorrecto con mayúsculas: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Transacción de Mayor Monto: ID 1 - 150.00")) {
		t.Errorf("Mayor monto incorrecto con mayúsculas: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Crédito: 1 Débito: 1")) {
		t.Errorf("Conteo incorrecto con mayúsculas: %s", output)
	}
}
