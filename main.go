package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// Transaccion que se extraerá del CSV
type Transaccion struct {
	ID    string
	Tipo  string
	Monto float64
}

func main() {
	// Verificar argumentos de línea de comandos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go [archivo_csv]")
		os.Exit(1)
	}

	// Obtener ruta del archivo CSV
	archivoCSV := os.Args[1]

	// Procesar el archivo CSV
	transacciones, err := leerTransaccionesDesdeCSV(archivoCSV)
	if err != nil {
		log.Fatalf("Error al leer el archivo CSV: %v", err)
	}

	// Generar reporte
	generarReporte(transacciones)
}

// leerTransaccionesDesdeCSV lee el archivo CSV y devuelve un slice de transacciones
func leerTransaccionesDesdeCSV(rutaArchivo string) ([]Transaccion, error) {
	// Abrir el archivo CSV
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer archivo.Close()

	// Crear un nuevo lector CSV
	reader := csv.NewReader(archivo)

	// Lee la primera línea (encabezados)
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error al leer los encabezados del CSV: %w", err)
	}

	// Slice para almacenar las transacciones
	var transacciones []Transaccion

	// Leer el resto de las líneas
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error al leer una línea del CSV: %w", err)
		}

		// Convertir el monto a float64
		monto, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("error al convertir el monto a número: %w", err)
		}

		// Crear y añadir la transacción
		transaccion := Transaccion{
			ID:    record[0],
			Tipo:  record[1],
			Monto: monto,
		}
		transacciones = append(transacciones, transaccion)
	}

	return transacciones, nil
}

// generarReporte procesa las transacciones y muestra el reporte en la terminal
func generarReporte(transacciones []Transaccion) {
	var balanceFinal float64
	var transaccionMayorMonto Transaccion
	conteoCredito := 0
	conteoDebito := 0

	// Procesar cada transacción
	for _, transaccion := range transacciones {
		// Verificar si es la transacción con mayor monto
		if transaccion.Monto > transaccionMayorMonto.Monto {
			transaccionMayorMonto = transaccion
		}

		// Actualizar contadores según el tipo de transacción
		switch strings.ToLower(transaccion.Tipo) {
		case "crédito", "credito":
			balanceFinal += transaccion.Monto
			conteoCredito++
		case "débito", "debito":
			balanceFinal -= transaccion.Monto
			conteoDebito++
		}
	}

	// Mostrar el reporte
	fmt.Println("Reporte de Transacciones")
	fmt.Println("---------------------------------------------")
	fmt.Printf("Balance Final: %.2f\n", balanceFinal)
	fmt.Printf("Transacción de Mayor Monto: ID %s - %.2f\n", transaccionMayorMonto.ID, transaccionMayorMonto.Monto)
	fmt.Printf("Conteo de Transacciones: Crédito: %d Débito: %d\n", conteoCredito, conteoDebito)
}
