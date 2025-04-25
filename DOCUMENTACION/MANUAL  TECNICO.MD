# Manual de Tecnico
**UNIVERSIDAD DE SAN CARLOS DE GUATEMALA**  
**FACULTAD DE INGENIERÍA**     
**CATEDRÁTICO:** ING. WILLIAM ESCOBAR
**TUTOR  ACADÉMICO:** KEVIN GARCÍA
**Nombre:** Brandon Antonio Marroquín Pérez  
**Carnet:** 202300813  
**Fecha de Entrega:** 30 de marzo de 2025  
**Curso:** lab MANEJO E IMPLEMENTACION DE ARCHIVOS
**Semestre:** SEXTO SEMESTRE 

---
## Introducción  
El proyecto consiste en desarrollar una aplicación web para la gestión de un sistema de archivos **EXT2**, simulando su comportamiento en un entorno distribuido. La solución integra un **frontend** moderno (con frameworks como Angular, React o Vue.js) y un **backend** en Go, que manejará las operaciones sobre archivos binarios (`.mia`) que emulan discos virtuales. El objetivo es permitir a los usuarios interactuar con el sistema de archivos mediante comandos, gestionar particiones, usuarios, permisos y generar reportes visuales con Graphviz.

---

## Objetivos  

### Objetivo General  
Desarrollar una aplicación web que simule un sistema de archivos **EXT2**, permitiendo la gestión de discos virtuales, particiones, usuarios, permisos y operaciones básicas de archivos, con una arquitectura cliente-servidor y reportes gráficos.  

### Objetivos Específicos  
1. **Implementar estructuras de datos clave** (MBR, EBR, inodos, bloques) en archivos binarios (`.mia`) para emular el sistema de archivos EXT2, respetando sus limitaciones y comportamientos teóricos.  
2. **Diseñar una interfaz web** que permita ejecutar comandos para administrar discos, particiones, usuarios y archivos, con áreas de entrada/salida y soporte para scripts (`.smia`).  

---

## Alcance del Sistema  
- **Frontend**: Interfaz web con áreas para comandos y resultados, carga de scripts y visualización de salidas.  
- **Backend**: APIs RESTful en Go para manejar operaciones sobre discos virtuales (creación, montaje, formateo) y estructuras EXT2 (inodos, bloques, bitmaps).  
- **Funcionalidades clave**:  
  - Comandos `MKDISK`, `FDISK`, `MOUNT`, `MKFS`, `LOGIN`, `MKFILE`, `MKDIR`.  
  - Gestión de usuarios/grupos (`MKUSR`, `RMGRP`, etc.) y permisos.  
  - Generación de reportes (`mbr`, `disk`, `tree`, etc.) con Graphviz.  
- **Limitaciones**:  
  - Solo soporta sistemas de archivos EXT2.  
  - No utiliza hardware físico (discos simulados en archivos binarios).  

---

## Consideraciones Técnicas  
1. **Lenguajes/Frameworks**:  
   - **Backend**: Go (obligatorio) para manejo eficiente de concurrencia y operaciones de bajo nivel.  
   - **Frontend**: Opcional (Angular recomendado).  
2. **Entorno de Ejecución**:  
   - Requiere una distribución GNU/Linux física para la calificación.  
3. **Restricciones**:  
   - Prohibido usar estructuras en memoria (listas, árboles) para gestionar archivos/carpetas.  
   - Los archivos `.mia` no deben cambiar de tamaño después de su creación.  
4. **Reportes**:  
   - Generados con Graphviz para validar funcionalidades (ej: `rep -name=mbr`).  

---

## Requisitos Recomendados  
- **Backend**:  
  - Go 1.20+ (para APIs RESTful y manejo de archivos binarios).  
  - Librerías: `os`, `encoding/binary`, `time`, `math`.  
- **Frontend**:  
  - Node.js + Angular/React/Vue.js.  
  - Librerías para gráficos: Graphviz (para reportes).  
- **Sistema Operativo**: Distribución Linux (Ubuntu 22.04+ recomendado).  

---

## Descripción de la Solución  
1. **Arquitectura**:  
   - **Cliente-Servidor**: Frontend se comunica con el backend mediante APIs REST.  
   - **Discos Virtuales**: Archivos `.mia` almacenan estructuras EXT2 (MBR, particiones, inodos).  
2. **Flujo de Operaciones**:  
   - El usuario ingresa comandos en el frontend (ej: `mkdisk -size=10 -path=/disco.mia`).  
   - El backend procesa el comando, actualiza el archivo binario y devuelve resultados.  
3. **Persistencia**:  
   - Todos los cambios se escriben directamente en los archivos `.mia` (sin caché en memoria).  
4. **Seguridad**:  
   - Permisos UGO (User-Group-Other) y usuario `root` con acceso total.  

---

## Conclusión  
Este proyecto permite entender a profundidad el funcionamiento de **EXT2** mediante una implementación práctica, combinando conceptos de sistemas de archivos, gestión de memoria y desarrollo full-stack. La solución fomenta el uso de buenas prácticas (arquitectura limpia, documentación) y valida su correcto funcionamiento mediante reportes gráficos. Además, sienta las bases para extensiones futuras (ej: soporte para otros sistemas de archivos).  


# Logica del programa
# Explicación Detallada de la Estructura y Conexión entre los Módulos Frontend y Backend

## 1. Estructura del Frontend (Angular)

El frontend está desarrollado con Angular, un framework de JavaScript que permite crear aplicaciones web interactivas. En este caso, se utiliza un servicio (`AnalyzerService`) para manejar las solicitudes HTTP hacia el backend.

### Explicación del `AnalyzerService`:

- **`HttpClient`**: Angular proporciona el `HttpClient` para interactuar con las APIs de backend. Este servicio permite enviar solicitudes HTTP (como `POST`) al backend.
- **`analyze(input: string)`**: Esta es la función que envía una solicitud al backend con el comando (`input`) proporcionado. El comando es un texto que el backend debe procesar y analizar.
- **Envío de datos**: En el método `analyze`, se utiliza `http.post()` para enviar el comando al backend. Los datos enviados incluyen un objeto con una propiedad `command`, que contiene el valor del comando a analizar.

```typescript
analyze(input: string): Observable<any> {
  console.log('Enviando al servidor:', { command: input });
  return this.http.post<any>(this.apiUrl, { command: input });
}


# Backend

El backend está construido con Go y el framework Fiber, que es ligero y optimizado para aplicaciones web.

## Explicación del código del backend

### Middleware CORS

Se utiliza `cors` para permitir que las solicitudes del frontend (en el navegador) puedan comunicarse con el servidor. Sin esto, las solicitudes podrían ser bloqueadas debido a restricciones de seguridad.

### Ruta `/analyze`

El backend escucha las solicitudes `POST` en la ruta `/analyze`. Cuando se recibe una solicitud, se procesa el comando enviado en el cuerpo de la solicitud.

### Procesamiento de comandos

- Se usa `BodyParser` para leer el cuerpo de la solicitud, que se espera sea un JSON con la propiedad `command`.
- El comando es procesado y se divide en varias líneas si contiene saltos de línea (`\n`).
- Se llama a la función `analyzer.Analyzer(cmd)` para procesar cada comando. Si hay un error al procesar un comando, se agrega un mensaje de error al resultado.

### Respuesta

Al finalizar, el backend responde al frontend con un objeto JSON que contiene un atributo `output`, que es el resultado del procesamiento de los comandos.

```go
app.Post("/analyze", func(c *fiber.Ctx) error {
  var req CommandRequest
  if err := c.BodyParser(&req); err != nil {
    return c.Status(400).JSON(CommandResponse{
      Output: "Error: Petición inválida",
    })
  }

  commands := strings.Split(req.Command, "\n")
  output := ""
  for _, cmd := range commands {
    if strings.TrimSpace(cmd) == "" {
      continue
    }
    result, err := analyzer.Analyzer(cmd)
    if err != nil {
      output += fmt.Sprintf("Error: %s\n", err.Error())
    } else {
      output += fmt.Sprintf("%s\n", result)
    }
  }

  if output == "" {
    output = "No se ejecutó ningún comando"
  }

  return c.JSON(CommandResponse{
    Output: output,
  })
})
```
Este bloque de código en Markdown está formateado para resaltar la sintaxis de Go, lo que facilita su lectura y comprensión.


## Conexión entre el Frontend y el Backend

El Frontend (Angular) y el Backend (Go) se comunican a través de HTTP. En este flujo:

1. El Frontend envía una solicitud `POST` al backend con el comando a analizar.
2. El Backend recibe la solicitud, procesa el comando y devuelve el resultado en formato JSON.
3. El Frontend recibe la respuesta y la presenta al usuario.

## Diagrama de Flujo

```plaintext
+---------------------+          HTTP POST          +-----------------------+
|                     | ----------------------->   |                       |
|    Frontend         |     /analyze (command)     |      Backend (Go)     |
|   (Angular)         |                           |                       |
|                     | <----------------------    |                       |
|    (User Input)     |   (command result as JSON) |                       |
|                     |                           |                       |
+---------------------+                           +-----------------------+

```
## Descripción del diagrama

**Frontend (Angular):**
El usuario ingresa un comando, y el frontend envía ese comando al backend mediante la API `/analyze`.

**Backend (Go):**
El backend procesa el comando, lo analiza y devuelve un resultado o mensaje de error.

**Frontend (Angular):**
El frontend recibe la respuesta y la muestra al usuario, quien puede interactuar con el sistema.

## Simulación de Sistema de Archivos EXT2

Aunque este código no implementa directamente un sistema de archivos, el backend (Go) podría estar procesando comandos relacionados con un sistema de archivos, como listar archivos, leer contenido, o manipular archivos. En el flujo de trabajo, el frontend (Angular) enviaría los comandos específicos del sistema de archivos (por ejemplo, `ls`, `cat`, `mkdir`, etc.) al backend para que los procese y devuelva el resultado.

## Explicación de las Estructuras de Datos

### MBR (Master Boot Record)
El MBR es la primera estructura en el disco y contiene información crucial sobre las particiones del disco.

- **Mbr_size**: Tamaño del MBR en bytes.
- **Mbr_creation_date**: Fecha y hora de creación del MBR.
- **Mbr_disk_signature**: Firma única del disco.
- **Mbr_disk_fit**: Tipo de ajuste (Best Fit, First Fit, Worst Fit).
- **Mbr_partitions**: Arreglo de particiones (hasta 4 particiones primarias).

El MBR se serializa y deserializa desde y hacia un archivo binario. También incluye métodos para obtener particiones disponibles, montar particiones y verificar la existencia de particiones.

### EBR (Extended Boot Record)
El EBR se utiliza para particiones lógicas y contiene información similar a la del MBR, pero específica para cada partición lógica.

- **Part_mount**: Indica si la partición está montada.
- **Part_fit**: Tipo de ajuste de la partición.
- **Part_start**: Byte de inicio de la partición.
- **Part_size**: Tamaño total de la partición en bytes.
- **Part_next**: Byte donde está el próximo EBR (-1 si no hay siguiente).
- **Part_name**: Nombre de la partición.

El EBR se serializa y deserializa desde y hacia un archivo binario. También incluye métodos para crear particiones lógicas y actualizar su estado.

### Inodos
Los inodos son estructuras que almacenan metadatos de archivos y directorios.

- **I_uid**: ID del usuario propietario.
- **I_gid**: ID del grupo propietario.
- **I_size**: Tamaño del archivo o directorio.
- **I_atime, I_ctime, I_mtime**: Tiempos de acceso, creación y modificación.
- **I_block**: Arreglo de bloques asociados.
- **I_type**: Tipo de inodo (archivo o directorio).
- **I_perm**: Permisos del inodo.

Los inodos se serializan y deserializan desde y hacia un archivo binario. También incluyen métodos para imprimir sus atributos.

### Bloques
Los bloques son unidades de almacenamiento que contienen datos de archivos o directorios.

- **FileBlock**: Bloque que contiene datos de un archivo.
- **FolderBlock**: Bloque que contiene datos de un directorio.

Ambos tipos de bloques se serializan y deserializan desde y hacia un archivo binario. También incluyen métodos para imprimir sus contenidos.

### SuperBloque
El SuperBloque contiene información global sobre el sistema de archivos.

- **S_filesystem_type**: Tipo de sistema de archivos.
- **S_inodes_count**: Número total de inodos.
- **S_blocks_count**: Número total de bloques.
- **S_free_inodes_count**: Número de inodos libres.
- **S_free_blocks_count**: Número de bloques libres.
- **S_mtime, S_umtime**: Tiempos de montaje y desmontaje.
- **S_mnt_count**: Número de veces que se ha montado el sistema de archivos.
- **S_magic**: Número mágico para verificación.
- **S_inode_size, S_block_size**: Tamaños de inodos y bloques.
- **S_first_ino, S_first_blo**: Primer inodo y bloque disponibles.
- **S_bm_inode_start, S_bm_block_start**: Inicio de los bitmaps de inodos y bloques.
- **S_inode_start, S_block_start**: Inicio de los inodos y bloques.

El SuperBloque se serializa y deserializa desde y hacia un archivo binario. También incluye métodos para imprimir sus atributos y generar diagramas de bloques.

### Bitmaps
Los bitmaps son estructuras que indican qué inodos y bloques están en uso o libres.

- **Bitmap de inodos**: Indica qué inodos están en uso.
- **Bitmap de bloques**: Indica qué bloques están en uso.

Los métodos `CreateBitMaps`, `UpdateBitmapInode` y `UpdateBitmapBlock` se encargan de crear y actualizar estos bitmaps en el archivo binario.

### Funciones Adicionales
- **CreateUsersFileExt2**: Crea el archivo `users.txt` en el sistema de archivos.
- **createFolderInInodeExt2**: Crea una carpeta en un inodo específico.
- **GenerateBlocksDot**: Genera un diagrama de bloques en formato DOT y lo convierte a PNG usando Graphviz.

## Organización y Gestión en el Archivo Binario `.mia`
El archivo binario `.mia` contiene todas las estructuras mencionadas (MBR, EBR, inodos, bloques, SuperBloque, bitmaps) organizadas de manera secuencial. Cada estructura tiene un offset específico dentro del archivo, lo que permite un acceso eficiente y una gestión adecuada del espacio en disco.

Este sistema de archivos simulado utiliza estas estructuras para gestionar archivos y directorios, permitiendo operaciones como crear, leer, escribir y eliminar archivos y directorios, así como montar y desmontar particiones.


## Descripción de los Comandos Implementados

### FDISK
El comando `FDISK` se utiliza para crear particiones en un disco.

- **Parámetros**:
  - `-size`: Tamaño de la partición (requerido).
  - `-unit`: Unidad de medida del tamaño (K, M, B) (opcional, por defecto M).
  - `-fit`: Tipo de ajuste (BF, FF, WF) (opcional, por defecto WF).
  - `-path`: Ruta del archivo del disco (requerido).
  - `-type`: Tipo de partición (P, E, L) (opcional, por defecto P).
  - `-name`: Nombre de la partición (requerido).

- **Ejemplo de uso**:
  ```sh
  
    fdisk -type=P -unit=B -name=PartErr -size=10485760 -path=/home/brandon/Calificacion_MIA/Discos/DiscoN.mia -fit=BF
## LOGIN
- El comando LOGIN se utiliza para iniciar sesión en una partición montada.

# Parámetros:

- user: Nombre de usuario (requerido).
- pass: Contraseña (requerida).
- id: ID del disco (requerido).
# Ejemplo de uso:

--Copiar
- login -user=root -pass=123 -id=062A3E2D
- Efectos:

# Verifica las credenciales del usuario.
Inicia sesión en la partición montada correspondiente al ID proporcionado.
## LOGOUT
El comando LOGOUT se utiliza para cerrar la sesión actual.

# Parámetros:

Ninguno.
Ejemplo de uso:

- Copiar
- logout
-Efectos:

- Cierra la sesión actual.
## MKDIR
- El comando MKDIR se utiliza para crear directorios en una partición montada.

- Parámetros:

- path: Ruta del directorio a crear (requerido).
- p: Opcional, crea directorios padres si no existen.
# Ejemplo de uso:

# Copiar
- mkdir -p -path=/home/user/docs/usac
- Efectos:

# Crea el directorio especificado en la partición montada.
- Si se usa la opción -p, crea los directorios padres necesarios.
## MKDISK
- El comando MKDISK se utiliza para crear un nuevo disco virtual.

# Parámetros:

-size: Tamaño del disco (requerido).
-unit: Unidad de medida del tamaño (K, M) (opcional, por defecto M).
-fit: Tipo de ajuste (BF, FF, WF) (opcional, por defecto FF).
-path: Ruta del archivo del disco (requerido).
# Ejemplo de uso:

# Copiar
- mkdisk -size=3000 -unit=K -path=/home/user/Disco1.mia
- Efectos:

# Crea un nuevo archivo de disco con el tamaño especificado.
- Inicializa el MBR del disco.
## MKFS
- El comando MKFS se utiliza para formatear una partición con un sistema de archivos específico.

## Parámetros:

- id: ID del disco (requerido).
- type: Tipo de formato (full) (opcional, por defecto full).
- fs: Sistema de archivos (2fs, 3fs) (opcional, por defecto 2fs).
# Ejemplo de uso:

# Copiar
- mkfs -id=vd1 -type=full -fs=2fs
- Efectos:

Formatea la partición con el sistema de archivos especificado.
Inicializa el superbloque y los bitmaps de inodos y bloques.

## MOUNT
El comando MOUNT se utiliza para montar una partición de un disco.

# Parámetros:

- path: Ruta del archivo del disco (requerido).
- name: Nombre de la partición (requerido).
Ejemplo de uso:

Copiar
mount -path=/home/Disco1.mia -name=Part1
Efectos:

Monta la partición especificada y actualiza el MBR para reflejar el estado montado.
## REP
El comando REP se utiliza para generar reportes del sistema de archivos.

- Parámetros:

- id: ID del disco (requerido).
- path: Ruta del archivo del disco (requerido).
- name: Nombre del reporte (mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls, tree) (requerido).
- path_file_ls: Ruta del archivo ls (opcional).
# Ejemplo de uso:

Copiar
rep -id=vd1 -path=/home/Disco1.mia -name=mbr
- Efectos:

Genera un reporte del tipo especificado para la partición montada.
## RMDISK
El comando RMDISK se utiliza para eliminar un disco virtual.

- Parámetros:

- path: Ruta del archivo del disco (requerido).
Ejemplo de uso:

- Copiar
- rmdisk -path="/home/mis discos/Disco4.mia"
- Efectos:

Elimina el archivo del disco especificado.
Copiar
Este formato Markdown es claro y estructurado, ideal para documentación técnica o para incluir en plataformas que soporten Markdown.

