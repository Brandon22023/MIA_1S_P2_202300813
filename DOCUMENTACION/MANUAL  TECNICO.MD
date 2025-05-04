# Manual de Tecnico
**UNIVERSIDAD DE SAN CARLOS DE GUATEMALA**  
**FACULTAD DE INGENIERÍA**     
**CATEDRÁTICO:** ING. WILLIAM ESCOBAR
**TUTOR  ACADÉMICO:** KEVIN GARCÍA
**Nombre:** Brandon Antonio Marroquín Pérez  
**Carnet:** 202300813  
**Fecha de Entrega:** 4 de mayo de 2025  
**Curso:** lab MANEJO E IMPLEMENTACION DE ARCHIVOS
**Semestre:** SEXTO SEMESTRE 

---
## Introducción  
El proyecto consiste en desarrollar una aplicación web para la gestión y simulación de sistemas de archivos **EXT2** y **EXT3**, permitiendo experimentar tanto con la estructura clásica de EXT2 como con el soporte de journaling de EXT3. La solución integra un **frontend** moderno (Angular, React o Vue.js) y un **backend** en Go, que maneja operaciones sobre archivos binarios (`.mia`) que emulan discos virtuales. El objetivo es que los usuarios puedan interactuar con ambos sistemas de archivos mediante comandos, gestionar particiones, usuarios, permisos y generar reportes visuales con Graphviz, explorando así las diferencias y ventajas de cada sistema.

---

### Objetivo General  
Desarrollar una aplicación web que simule sistemas de archivos **EXT2** y **EXT3**, permitiendo la gestión de discos virtuales, particiones, usuarios, permisos y operaciones básicas de archivos, con una arquitectura cliente-servidor y reportes gráficos, mostrando las ventajas del journaling en EXT3.

### Objetivos Específicos  
1. **Implementar estructuras de datos clave** (MBR, EBR, inodos, bloques, journal) en archivos binarios (`.mia`) para emular los sistemas de archivos EXT2 y EXT3, respetando sus limitaciones y comportamientos teóricos.  
2. **Diseñar una interfaz web** que permita ejecutar comandos para administrar discos, particiones, usuarios y archivos, con áreas de entrada/salida y soporte para scripts (`.smia`).  
3. **Comparar y demostrar** el funcionamiento y ventajas de EXT3 respecto a EXT2, especialmente en cuanto a integridad y recuperación ante fallos mediante journaling.

---

## Alcance del Sistema  
- **Frontend**: Interfaz web con áreas para comandos y resultados, carga de scripts y visualización de salidas.  
- **Backend**: APIs RESTful en Go para manejar operaciones sobre discos virtuales (creación, montaje, formateo) y estructuras EXT2/EXT3 (inodos, bloques, bitmaps, journal).  
- **Funcionalidades clave**:  
  - Comandos `MKDISK`, `FDISK`, `MOUNT`, `MKFS`, `LOGIN`, `MKFILE`, `MKDIR`, `MKUSR`, `RMGRP`, `COPY`, `MOVE`, `REP`, `RMDISK`, `LOGOUT`, `UNMOUNT`.
  - Comandos específicos de EXT3: `JOURNAL`, `RECOVERY`, `SHOWJOURNAL`.
  - Gestión de usuarios/grupos y permisos.
  - Generación de reportes (`mbr`, `disk`, `tree`, etc.) con Graphviz.  
- **Limitaciones**:  
  - Solo soporta sistemas de archivos EXT2 y EXT3 simulados.  
  - No utiliza hardware físico (discos simulados en archivos binarios).  



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
Este proyecto permite entender a profundidad el funcionamiento de los sistemas de archivos **EXT2** y **EXT3** mediante una implementación práctica, combinando conceptos de sistemas de archivos, gestión de memoria, journaling y desarrollo full-stack. La solución fomenta el uso de buenas prácticas (arquitectura limpia, documentación) y valida su correcto funcionamiento mediante reportes gráficos. Además, sienta las bases para extensiones futuras y permite comparar la robustez y ventajas de EXT3 frente a EXT2, especialmente en escenarios de recuperación ante fallos.


| Comando      | Descripción                                                                                   | Parámetros principales                                                                                  | Soporte EXT2 | Soporte EXT3 |
|--------------|----------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|--------------|--------------|
| MKDISK       | Crea un nuevo disco virtual.                                                                  | -size, -unit, -fit, -path                                                                              | ✔️           | ✔️           |
| RMDISK       | Elimina un disco virtual.                                                                     | -path                                                                                                  | ✔️           | ✔️           |
| FDISK        | Crea particiones en un disco.                                                                 | -size, -unit, -fit, -path, -type, -name                                                                | ✔️           | ✔️           |
| MOUNT        | Monta una partición de un disco.                                                              | -path, -name                                                                                           | ✔️           | ✔️           |
| UNMOUNT      | Desmonta una partición montada.                                                               | -id                                                                                                    | ✔️           | ✔️           |
| MKFS         | Formatea una partición con EXT2 o EXT3.                                                       | -id, -type, -fs                                                                                        | ✔️           | ✔️           |
| LOGIN        | Inicia sesión en una partición montada.                                                       | -user, -pass, -id                                                                                      | ✔️           | ✔️           |
| LOGOUT       | Cierra la sesión actual.                                                                      | Ninguno                                                                                                | ✔️           | ✔️           |
| MKDIR        | Crea directorios en una partición montada.                                                    | -path, -p                                                                                              | ✔️           | ✔️           |
| MKFILE       | Crea archivos en una partición montada.                                                       | -path, -size, -cont                                                                                    | ✔️           | ✔️           |
| COPY         | Copia archivos o carpetas dentro del sistema de archivos.                                     | -path, -destino                                                                                        | ✔️           | ✔️           |
| MOVE         | Mueve archivos o carpetas dentro del sistema de archivos.                                     | -path, -destino                                                                                        | ✔️           | ✔️           |
| REMOVE       | Elimina archivos del sistema de archivos.                                                     | -path                                                                                                  | ✔️           | ✔️           |
| RENAME       | Cambia el nombre de un archivo o carpeta.                                                     | -path, -name                                                                                           | ✔️           | ✔️           |
| FIND         | Busca archivos o carpetas por patrón en el sistema de archivos.                               | -path, -name                                                                                           | ✔️           | ✔️           |
| REP          | Genera reportes del sistema de archivos (mbr, disk, inode, block, bm_inode, bm_block, sb, etc.)| -id, -path, -name, -path_file_ls                                                                       | ✔️           | ✔️           |
| JOURNAL      | Muestra el contenido del journal (solo EXT3).                                                 | -id                                                                                                    | ❌           | ✔️           |

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

### 1. analyze

El endpoint `/analyze` recibe comandos como texto, los separa por líneas y los procesa uno a uno usando el analizador principal. Entre los comandos soportados se encuentran los nuevos: `COPY`, `MOVE`, `REMOVE`, `RENAME`, `FIND`, además de los clásicos como `MKDISK`, `RMDISK`, `MOUNT`, etc.

**Fragmento relevante:**
```go
app.Post("/analyze", func(c *fiber.Ctx) error {
    ...
    for _, cmd := range commandsList {
        ...
        result, err := analyzer.Analyzer(cmd)
        ...
    }
    ...
})
```
**Descripción:**  
Cada comando es analizado y ejecutado, permitiendo manipular archivos y carpetas (copiar, mover, eliminar, renombrar, buscar) recorriendo inodos y bloques para mantener la persistencia y consistencia en el archivo `.mia`.

---

### 2. Extracción y Exportación de Información

Después de ejecutar comandos, el backend exporta información de los discos y extrae archivos `.txt` válidos de la partición activa, actualizando las listas globales.

**Fragmento relevante:**
```go
// Exportar información de cada path
for _, diskPath := range paths {
    err := commands.ExportDiskInfo(diskPath)
    ...
}
// Extraer archivos .txt válidos
partitionID, err := stores.GetActivePartitionID()
if err == nil {
    sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(partitionID)
    if err == nil {
        sb.ExtractTxtFiles(partitionPath, partitionID)
    }
}
```
**Descripción:**  
Esto permite que el frontend consulte información actualizada sobre discos, particiones, carpetas y archivos `.txt` tras cada operación.

---

### 3. Nuevos Endpoints para Carpetas y Archivos

Se agregaron endpoints para consultar carpetas y archivos `.txt` válidos en la partición activa:

- `/folders`: Devuelve las rutas de carpetas válidas.
- `/txtfiles`: Devuelve los archivos `.txt` válidos.

**Fragmento relevante:**
```go
app.Get("/folders", func(c *fiber.Ctx) error {
    ...
    for _, path := range global.ValidPaths {
        ...
    }
    ...
})

app.Get("/txtfiles", func(c *fiber.Ctx) error {
    ...
    txtFiles, err := sb.GetTxtFiles(partitionPath, partitionID)
    ...
})
```
**Descripción:**  
Estos endpoints permiten al frontend mostrar la estructura de carpetas y archivos disponibles para operaciones como copiar, mover, eliminar o buscar.



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

## Uso de AWS, Docker y Termius para Despliegue en la Nube

Para facilitar el acceso, la disponibilidad y la portabilidad del sistema, se utilizarán servicios de Amazon Web Services (AWS) junto con herramientas modernas de despliegue y administración:

- **Amazon EC2**: Se utilizará para alojar y ejecutar tanto el backend como el frontend de la aplicación, permitiendo que los usuarios accedan desde cualquier lugar.
- **Amazon S3**: Se empleará para el almacenamiento de archivos y respaldos importantes, asegurando la persistencia y disponibilidad de los datos generados por el sistema.
- **Docker**: El backend será empaquetado y desplegado usando Docker, lo que facilita la portabilidad, el despliegue y la gestión de dependencias, asegurando que el entorno de ejecución sea consistente en cualquier servidor.
- **Termius**: Se utilizará Termius como cliente SSH para subir y administrar tanto el backend (contenedor Docker) como el frontend estático en la instancia EC2, permitiendo una gestión remota eficiente y segura.

Esta integración permite escalar la solución, mejorar la seguridad, facilitar la administración de los recursos y garantizar que el sistema pueda ser desplegado y actualizado de manera sencilla y confiable en la nube.

## Simulación de Sistema de Archivos EXT2

Aunque este código no implementa directamente un sistema de archivos, el backend (Go) podría estar procesando comandos relacionados con un sistema de archivos, como listar archivos, leer contenido, o manipular archivos. En el flujo de trabajo, el frontend (Angular) enviaría los comandos específicos del sistema de archivos (por ejemplo, `ls`, `cat`, `mkdir`, etc.) al backend para que los procese y devuelva el resultado.

## EXT3

El sistema de archivos **EXT3** es una extensión de EXT2 que añade soporte para journaling, lo que permite una mayor integridad y recuperación ante fallos. En este proyecto, la implementación de EXT3 se basa en las mismas estructuras fundamentales de EXT2 (MBR, EBR, inodos, bloques, superbloque), pero incorpora la estructura adicional de **Journal** para registrar operaciones críticas.

### Características principales de EXT3 en la implementación

- **Journaling**: Cada operación importante (como la creación de archivos o carpetas) se registra en una estructura de journal antes de aplicarse al sistema de archivos. Esto permite recuperar el sistema en caso de un fallo inesperado.
- **Estructura Journal**: Se implementa la estructura `Journal`, que almacena información sobre la operación realizada, la ruta afectada, el contenido y la fecha/hora.
- **Compatibilidad**: EXT3 mantiene compatibilidad con las estructuras de EXT2, permitiendo migrar o montar particiones EXT2 como EXT3.
- **Persistencia**: El journal se almacena directamente en el archivo binario `.mia`, siguiendo la restricción de no usar estructuras en memoria para la gestión de archivos/carpetas.

### Estructuras adicionales

- **Journal**
  - `J_count`: Número de operación registrada.
  - `J_content`: Información de la operación (tipo, ruta, contenido, fecha).

### Ejemplo de flujo de journaling

1. El usuario ejecuta una operación (por ejemplo, crear una carpeta).
2. Antes de modificar los inodos o bloques, se registra la operación en el journal.
3. Si la operación se completa correctamente, el sistema continúa normalmente.
4. En caso de fallo, el sistema puede consultar el journal para recuperar o revertir operaciones pendientes.

### Funciones clave implementadas

- `CreateUsersFileExt3`: Crea el archivo `users.txt` y registra la operación en el journal.
- `createFolderInInodeExt3`: Crea carpetas y registra cada operación en el journal antes de modificar las estructuras principales.
- Serialización y deserialización del journal para asegurar persistencia en disco.

### Ventajas de EXT3 sobre EXT2

- **Recuperación rápida**: Gracias al journaling, el sistema puede recuperarse rápidamente tras un fallo, sin necesidad de revisar todo el sistema de archivos.
- **Integridad**: Menor riesgo de corrupción de datos ante apagados inesperados o errores del sistema.

### Limitaciones

- El journaling implementado es básico y está orientado a operaciones de creación de archivos y carpetas.
- El rendimiento puede ser ligeramente inferior al de EXT2 debido al registro adicional de operaciones.

### Ejemplo de estructura Journal en Go

```go
type Journal struct {
    J_count   int32
    J_content Information
}

type Information struct {
    I_operation [10]byte
    I_path      [32]byte
    I_content   [64]byte
    I_date      float32
}
```

### Resumen

La implementación de EXT3 en este proyecto permite simular un sistema de archivos robusto y tolerante a fallos, manteniendo la compatibilidad con EXT2 y cumpliendo con los principios de persistencia y seguridad exigidos por el curso.
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


## Comandos Nuevos: Resumen y Ejemplo de Código

A continuación se describen los comandos nuevos implementados en el sistema, junto con un resumen funcional y un fragmento representativo de su código para entender su propósito y funcionamiento.

---

### COPY

**Propósito:**  
Copia archivos o carpetas (y su contenido) dentro del sistema de archivos, duplicando la estructura y los datos.

**Fragmento de código:**
```go
// ParseCopy procesa el comando copy
func ParseCopy(args []string) (string, error) {
    ...
    err = copyFolderRecursive(sb, partitionPath, srcPath, dstPath, true)
    ...
}
```
**Descripción:**  
`ParseCopy` recibe los paths de origen y destino, y llama a `copyFolderRecursive`, que recorre inodos y bloques para copiar carpetas y archivos, asegurando que la estructura y los datos se dupliquen correctamente en el destino.

---

### MOVE

**Propósito:**  
Mueve archivos o carpetas a otra ubicación, eliminando el original y actualizando las rutas en las estructuras internas.

**Fragmento de código:**
```go
// ParseMove procesa el comando move
func ParseMove(args []string) (string, error) {
    ...
    // Actualizar la lista global: reemplazar srcPath por dstFilePath
    ...
}
```
**Descripción:**  
`ParseMove` valida los paths, actualiza las rutas en la lista global y mueve el archivo/carpeta, asegurando que los datos y referencias se actualicen correctamente, eliminando el original tras la copia.

---

### REMOVE

**Propósito:**  
Elimina archivos del sistema de archivos, liberando los inodos y bloques asociados y actualizando las estructuras de carpetas.

**Fragmento de código:**
```go
func ParseRemove(args []string) (string, error) {
    ...
    err = RemoveFile(sb, partitionPath, tryPath)
    ...
}
```
**Descripción:**  
`ParseRemove` busca el archivo a eliminar recorriendo los inodos y bloques, y llama a `RemoveFile`, que elimina la referencia en la carpeta padre y libera los recursos asociados.

---

### RENAME

**Propósito:**  
Cambia el nombre de un archivo o carpeta, modificando la entrada correspondiente en el bloque de la carpeta padre.

**Fragmento de código:**
```go
func ParseRename(args []string) (string, error) {
    ...
    err = RenameFile(sb, partitionPath, tryPath, newName)
    ...
}
```
**Descripción:**  
`ParseRename` localiza el archivo/carpeta y llama a `RenameFile`, que actualiza el nombre en el bloque de la carpeta padre, permitiendo cambiar nombres sin alterar el contenido ni la ubicación.

---

### FIND

**Propósito:**  
Busca archivos o carpetas que coincidan con un patrón, mostrando la estructura encontrada en formato de árbol.

**Fragmento de código:**
```go
func ParseFind(args []string) (string, error) {
    ...
    tree := buildTree(matches, path)
    ...
}
```
**Descripción:**  
`ParseFind` filtra las rutas que coinciden con el patrón usando `matchPattern` y construye una representación tipo árbol con `buildTree`, facilitando la localización visual de archivos y carpetas.

---

> **Nota:** Todos estos comandos operan recorriendo inodos y bloques para manipular archivos y carpetas, asegurando la persistencia y consistencia en el archivo binario `.mia`.