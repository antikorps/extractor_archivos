# Extractor de archivos

Busca dentro de un directorio (incluyendo las subcarpetas) todos aquellos archivos con una extensión determinada para moverlos automáticamente a otra carpeta.

Su uso original se pensó para archivos multimedia, pero es completamente configurable. Solamente es necesario crear un archivo llamado extensiones.txt en la misma carpeta que el ejecutable e incorporar las extensiones deseadas (una en cada linea).

## Descargas

No es necesaria ningún tipo de instalación, únicamente ejecutar los binarios descargados.

* [GNU/Linux amd64](https://github.com/antikorps/extractor_archivos/raw/main/bin/extractor_archivos)
* [Windows amd64](https://github.com/antikorps/extractor_archivos/raw/main/bin/extractor_archivos.exe)

## Funcionamiento

Asegurar que el archivo tiene permisos de ejecución. En GNU/Linux tal vez sea necesario otorgarlo. La forma más rápida es mediante el siguiente comando:

```bash
sudo chmod +x 
```

Aunque también puede hacerse con la mayoría de exploradores de archivos (Botón derecho sobre el archivo > Propiedades > Permisos > Permitir ejecutar el archivo).

La primera pantalla es la descripción del programa y sus detalles.

![Descripción y detalles](https://i.imgur.com/Upx0fWn.png "Descripción y detalles")

A continuación aparecerá una lista con una colección de extensiones para seleccionar aquellas que se quieren utilizar. Recuerda que si necesitas una lista de extensiones personalizadas solo tienes que crear un archivo llamado extensiones.txt junto al ejecutable y escribir una extensión en cada línea.

Puedes seleccionar una, varias o todas. Para que sea más cómodo, en el caso de que se no se seleccione ninguna se entenderá que se quieren seleccionar todas.

![Selección de extensiones](https://i.imgur.com/rgDyPHh.png "Selección de extensiones")

El siguiente paso será seleccionar la carpeta/directorio/ruta que se quiere analizar. Se recorrerán su contenido (subcarpetas y archivos) buscando los archivos que coincidan con alguna de las extensiones válidas selecionadas.

![Directorio a analizar](https://i.imgur.com/tRTfTVH.png "Directorio a analizar")

Posteriormente aparecerá un resumen con todos los archivos encontrados antes de moverse. En este paso se pide la confirmación de los resultados, puede seleccionarse un archivo, varios o todos. Al igual que antes, en el caso de que no se seleccione ningún archivo se entenderá que se han seleccionado todos.

![Resumen de archivos](https://i.imgur.com/j8pRqNk.png "Resumen de archivos")

A continuación volverá a aparecer un selector de archivos en el que deberá escogerse el directorio/ruta/carpeta en el que se moverán todos los archivos encontrados.

![Directorio destino](https://i.imgur.com/atzwM5X.png "Directorio destino")

Tras las operaciones aparecerá una ventana con los resultados donde se mostrará el número de archivos movidos correctamente y el de los posibles errores. En el caso de que sea necesario se puede generar un informe de estos resultados en formato CSV, es opcional y puede ordenarse pulsando el botón "Sí, generar informe"".

![Resultados](https://i.imgur.com/MBUuheF.png "Resultados")

En el caso que haya sido necesario generar el informe aparecerá una nueva ventana confirmando su preparación y la ruta del archivo. Es fácil de recordar porque siempre se generará en la carpeta de destino, junto a los archivos movidos, con el nombre informe_resultados.csv

![Informe](https://i.imgur.com/8XOVOIA.png "Informe")

El informe en CSV contendrá la siguiente información:

* error
* error_mensaje
* nombre_archivo
* ruta_destino
* ruta_origen
