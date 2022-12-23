package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

type Archivo struct {
	Nombre string
	Ruta   string
}

type ResultadoOperacion struct {
	Error        string
	ErrorMensaje string
	Nombre       string
	RutaOrigen   string
	RutaDestino  string
}

type ManejadorArchivos struct {
	ArchivosEncontrados      []Archivo
	ArchivosMover            []string
	DirectorioBusqueda       string
	DirectorioDestino        string
	ExtensionesSeleccionadas []string
	ExtensionesValidas       []string
	GenerarInforme           bool
	Resultados               []ResultadoOperacion
	ResultadosErrores        int
	ResultadosExitos         int
}

func mostrarVentanaError(mensajeError string, critico bool) {
	zenity.Error(mensajeError,
		zenity.Title("Error"),
		zenity.ErrorIcon)
	if critico {
		log.Fatalln(mensajeError)
	}
}

func obtenerExtensionesArchivo(rutaExtensiones string) []string {
	var extensionesArchivo []string

	archivoExtensiones, archivoExtensionesError := os.Open(rutaExtensiones)
	if archivoExtensionesError != nil {
		return extensionesArchivo
	}
	escaner := bufio.NewScanner(archivoExtensiones)
	escaner.Split(bufio.ScanLines)

	for escaner.Scan() {
		extension := escaner.Text()
		if extension == "" {
			continue
		}
		extension = strings.TrimSpace(extension)
		if !strings.HasPrefix(extension, ".") {
			extension = "." + extension
		}
		extensionesArchivo = append(extensionesArchivo, extension)
	}

	return extensionesArchivo
}

func crearManejador() ManejadorArchivos {
	var extensiones []string

	ejecutable, ejecutableError := os.Executable()
	if ejecutableError != nil {
		log.Fatalln(ejecutableError)
	}
	rutaEjecutable := filepath.Dir(ejecutable)
	rutaArchivoExtensiones := filepath.Join(rutaEjecutable, "extensiones.txt")

	_, archivoExtensionesError := os.Stat(rutaArchivoExtensiones)
	if archivoExtensionesError == nil {
		extensiones = obtenerExtensionesArchivo(rutaArchivoExtensiones)
	}
	if len(extensiones) == 0 {
		extensiones = []string{".avi", ".mkv", ".mov", ".mp3", ".mp4", ".mpeg", ".mpg", ".ogg", ".ogv", ".qt ", ".rm", ".webm"}
	}

	return ManejadorArchivos{ExtensionesValidas: extensiones}
}

func (m *ManejadorArchivos) configurar() {
	inicioError := zenity.Info("Extrae los archivos por su tipo (extensión) de un directorio y subcarpetas.\n\n1- Selecciona el directorio que quieres analizar\n2- Escoge las extensiones de los archivos que buscas\n3- Confirma el resumen\n4- Selecciona la carpeta de destino a donde moverlos\n5- Genera un informe en CSV si lo necesitas\n\nSi quieres personalizar las extensiones:\n\na) crea un archivo llamado extensiones.txt en el mismo directorio que el ejecutable\nb) escribe cada extensión en una línea.",
		zenity.Title("Extractor de archivos"),
		zenity.Width(400),
		zenity.OKLabel("Continuar"),
		zenity.InfoIcon)
	if inicioError != nil {
		log.Fatalln(inicioError)
	}

	extensionesArchivos, extensionesArchivosError := zenity.ListMultiple(
		"Selecciona las extensiones de los archivos que quieres mover\nSi no seleccionas ninguna se procesarán todas",
		m.ExtensionesValidas,
		zenity.Title("Extensiones válidas"),
		zenity.OKLabel("Continuar"),
		zenity.CancelLabel("Salir"),
	)

	if extensionesArchivosError != nil {
		if extensionesArchivosError.Error() == "dialog canceled" {
			os.Exit(0)
		}
		mostrarVentanaError(extensionesArchivosError.Error(), true)
	}

	m.ExtensionesSeleccionadas = extensionesArchivos

	if len(extensionesArchivos) == 0 {
		m.ExtensionesSeleccionadas = m.ExtensionesValidas
	}

	carpetaOrigen, carpetaOrigenError := zenity.SelectFile(
		zenity.Title("Selecciona el directorio que quieres analizar"),
		zenity.Directory(),
		zenity.OKLabel("Seleccionar"),
		zenity.CancelLabel("Salir"))
	if carpetaOrigenError != nil {
		if carpetaOrigenError.Error() == "dialog canceled" {
			os.Exit(0)
		}
		mostrarVentanaError(carpetaOrigenError.Error(), true)
	}
	if carpetaOrigen == "" {
		os.Exit(0)
	}
	m.DirectorioBusqueda = carpetaOrigen

}

func comprobarExtensionValida(extensionesValidas []string, rutaArchivo string) bool {
	extensionArchivo := filepath.Ext(rutaArchivo)
	for _, v := range extensionesValidas {
		if extensionArchivo == v {
			return true
		}
	}
	return false
}

func (m *ManejadorArchivos) buscarArchivos() {
	var archivosEncontrados []Archivo

	busquedaArchivosError := filepath.Walk(m.DirectorioBusqueda,
		func(ruta string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			var archivo Archivo
			if comprobarExtensionValida(m.ExtensionesSeleccionadas, ruta) {
				archivo.Nombre = info.Name()
				archivo.Ruta = ruta
				archivosEncontrados = append(archivosEncontrados, archivo)
			}
			return nil
		})
	if busquedaArchivosError != nil {
		mostrarVentanaError(busquedaArchivosError.Error(), true)
	}
	m.ArchivosEncontrados = archivosEncontrados
}

func (m *ManejadorArchivos) confirmar() {
	var archivosMover []string

	var archivosEncontrados []string
	for _, v := range m.ArchivosEncontrados {
		archivosEncontrados = append(archivosEncontrados, v.Nombre)
	}
	mensajeListaArchivos := fmt.Sprint("Se han encontrado ", len(archivosEncontrados), " archivos.\nSelecciona los que quieres mover\nSi no seleccionas ninguno se moverán todos.\nEn el siguiente paso deberás seleccionar la carpeta a la que se moverán")
	archivosSeleccionados, archivosSeleccionadosError := zenity.ListMultiple(
		mensajeListaArchivos,
		archivosEncontrados,
		zenity.Title("Archivos encontrados"),
		zenity.OKLabel("Continuar"),
		zenity.CancelLabel("Salir"),
	)

	if archivosSeleccionadosError != nil {
		if archivosSeleccionadosError.Error() == "dialog canceled" {
			os.Exit(0)
		}
		mostrarVentanaError(archivosSeleccionadosError.Error(), true)
	}

	archivosMover = append(archivosMover, archivosSeleccionados...)
	if len(archivosMover) == 0 {
		for _, v := range m.ArchivosEncontrados {
			archivosMover = append(archivosMover, v.Nombre)
		}
	}

	m.ArchivosMover = archivosMover
}

func (m *ManejadorArchivos) seleccionarDestino() {
	carpetaDestino, carpetaDestinoError := zenity.SelectFile(
		zenity.Title("Selecciona el directorio al que mover los archivos"),
		zenity.Directory())
	if carpetaDestinoError != nil {
		if carpetaDestinoError.Error() == "dialog canceled" {
			os.Exit(0)
		}
		mostrarVentanaError(carpetaDestinoError.Error(), true)
	}
	if carpetaDestino == "" {
		os.Exit(0)
	}
	m.DirectorioDestino = carpetaDestino
}

func (m *ManejadorArchivos) moverArchivos() {
	for _, v := range m.ArchivosMover {
	compararArchivoConNombre:
		for _, archivo := range m.ArchivosEncontrados {
			if archivo.Nombre == v {
				rutaNuevoArchivo := filepath.Join(m.DirectorioDestino, archivo.Nombre)
				moverError := os.Rename(archivo.Ruta, rutaNuevoArchivo)

				var resultadoOperacion ResultadoOperacion
				resultadoOperacion.Nombre = archivo.Nombre
				resultadoOperacion.RutaOrigen = archivo.Ruta
				if moverError != nil {
					resultadoOperacion.Error = "sí"
					resultadoOperacion.ErrorMensaje = moverError.Error()
					resultadoOperacion.RutaDestino = ""
					m.ResultadosErrores++
				} else {
					resultadoOperacion.RutaDestino = rutaNuevoArchivo
					m.ResultadosExitos++
				}
				m.Resultados = append(m.Resultados, resultadoOperacion)
				break compararArchivoConNombre
			}
		}
	}
}

func (m *ManejadorArchivos) mostrarResumen() {
	mensajeResumen := fmt.Sprint("Se han movido correctamente ", m.ResultadosExitos, " archivos.\nSe ha producido ", m.ResultadosErrores, " errores.\n¿Quieres generar un informe en .csv?")

	resumenError := zenity.Question(mensajeResumen,
		zenity.Title("Resumen"),
		zenity.OKLabel("Sí, generar informe"),
		zenity.CancelLabel("No, Salir"),
		zenity.Width(400),
		zenity.QuestionIcon)

	if resumenError != nil {
		if resumenError.Error() == "dialog canceled" {
			os.Exit(0)
			return
		}
		mostrarVentanaError(resumenError.Error(), true)
	}
	m.GenerarInforme = true
}

func (m *ManejadorArchivos) generarInforme() {
	if !m.GenerarInforme {
		return
	}

	var filasCSV [][]string

	encabezado := []string{"error", "error_mensaje", "nombre_archivo", "ruta_origen", "ruta_destino"}
	filasCSV = append(filasCSV, encabezado)

	for _, v := range m.Resultados {
		linea := []string{v.Error, v.ErrorMensaje, v.Nombre, v.RutaOrigen, v.RutaDestino}
		filasCSV = append(filasCSV, linea)
	}

	rutaInforme := filepath.Join(m.DirectorioDestino, "informe_resultados.csv")
	informe, informeError := os.Create(rutaInforme)
	if informeError != nil {
		mostrarVentanaError(informeError.Error(), true)
	}
	defer informe.Close()

	manejadorCSV := csv.NewWriter(informe)
	escribirCSVError := manejadorCSV.WriteAll(filasCSV)
	if escribirCSVError != nil {
		mostrarVentanaError(escribirCSVError.Error(), true)
	}
	manejadorCSV.Flush()

	mensajeInforme := fmt.Sprint("Informe generado. Puedes consultarlo en: ", rutaInforme)
	zenity.Info(mensajeInforme,
		zenity.Title("Informe actualizado"),
		zenity.Width(400),
		zenity.InfoIcon)
}

func main() {
	manejador := crearManejador()
	manejador.configurar()
	manejador.buscarArchivos()
	manejador.confirmar()
	manejador.seleccionarDestino()
	manejador.moverArchivos()
	manejador.mostrarResumen()
	manejador.generarInforme()
}
