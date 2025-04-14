# Procesador de Transacciones Bancarias (CLI)

## Introducción
Este proyecto es una aplicación de línea de comandos (CLI) desarrollada en Go que procesa un archivo CSV con transacciones bancarias y genera un reporte con estadísticas clave, incluyendo el balance final, la transacción de mayor monto y el conteo de transacciones por tipo.

## Instrucciones de Ejecución

### Requisitos Previos
- Go 1.16 o superior instalado en el sistema

### Instalación
1. Clona este repositorio:
   ```
   git clone https://github.com/imjowend/interbank-academy-25
   cd interbank-academy-25
   ```

2. Compilar la aplicación:
   ```
   go build -o prueba-cli
   ```

### Ejecución
Para ejecutar la aplicación, usa el siguiente comando:

```
./prueba-cli transacciones.csv
```

Donde `transacciones.csv` es la ruta al archivo CSV que contiene las transacciones bancarias.

### Ejemplo de uso
1. Crea un archivo CSV con el siguiente contenido:
   ```
   id,tipo,monto
   1,Crédito,100.00
   2,Débito,50.00
   3,Crédito,200.00
   4,Débito,75.00
   5,Crédito,150.00
   ```

2. Ejecuta la aplicación:
   ```
   ./prueba-cli transacciones.csv
   ```

3. La salida en la terminal será:
   ```
   Reporte de Transacciones
   ---------------------------------------------
   Balance Final: 325.00
   Transacción de Mayor Monto: ID 3 - 200.00
   Conteo de Transacciones: Crédito: 3 Débito: 2
   ```

## Enfoque y Solución

La solución implementada sigue un enfoque modular y directo para procesar las transacciones bancarias:

1. **Lectura de datos**: La aplicación lee el archivo CSV línea por línea, extrayendo los datos de cada transacción.
2. **Procesamiento**: Cada transacción se analiza para:
    - Sumar o restar al balance final según su tipo
    - Identificar la transacción con el mayor monto
    - Contar el número de transacciones de cada tipo
3. **Generación del reporte**: La aplicación formatea y muestra los resultados en la terminal.

El código está diseñado para ser robusto, manejando posibles errores durante la lectura del archivo y la conversión de datos, lo que garantiza una experiencia de usuario sin problemas.

## Estructura del Proyecto

- `main.go`: Contiene la lógica principal de la aplicación, incluyendo:
    - Función `main()`: Punto de entrada de la aplicación
    - Función `leerTransaccionesDesdeCSV()`: Lee y procesa el archivo CSV
    - Función `generarReporte()`: Analiza las transacciones y genera el reporte final
- `README.md`: Este archivo con documentación detallada del proyecto
- `go.mod`: Archivo de configuración del módulo Go