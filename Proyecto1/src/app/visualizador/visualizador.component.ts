import { Component, EventEmitter, Output } from '@angular/core';
import { AnalyzerService } from '../analyzer.service';
import { CommonModule } from '@angular/common';
@Component({
  selector: 'app-visualizador',
  imports: [CommonModule],
  templateUrl: './visualizador.component.html',
  styleUrl: './visualizador.component.css'
})
export class VisualizadorComponent {

  @Output() volverEvent = new EventEmitter<void>(); // Evento para volver al componente principal
  disks: { name: string; size: string; fit: string; mounted_partitions: string }[] = []; // Lista de discos
  partitions: { 
    name: string; 
    size: string; 
    type: string; 
    start: string; 
    fit: string; 
    state: string; 
    id: string; 
  }[] = []; // Lista de particiones
  carpetas: { path: string; id: string; permissions: string; }[] = [];
  txt: { path: string; id: string; permissions: string; contenido: string; size: number }[] = [];
  currentPath: string = ''; // Ruta actual
  ruta: string[] = []; // Ruta desglosada para mostrar en la barra
  selectedDisk: string | null = null;
  selectedpart: string | null = null;
  selectedFileContent: string | null = null; // Contenido del archivo seleccionado
  constructor(private analyzerService: AnalyzerService) {}
  selectedpartId: string | null = null;
  
  ngOnInit(): void {
    this.loadDisks();
    this.loadFolders(); // Cargar carpetas desde el backend
    // Archivo .txt quemado para pruebas
  }
  
  loadFolders(): void {
    this.analyzerService.getFolders().subscribe({
      next: (response) => {
        // Agregar la propiedad "permissions" con un valor predeterminado
        this.carpetas = response.map((folder) => ({
          ...folder,
          permissions: '664', // Asignar permisos predeterminados
        }));
      },
      error: (err) => {
        console.error('Error al cargar las carpetas:', err);
      }
    });
  }
  loadTxtFiles(): void {
  this.analyzerService.getTxtFiles().subscribe({
    next: (response) => {
      console.log('Respuesta real del backend txtfiles:', response);
      this.txt = response.map(file => ({
        ...file,
        permissions: '664',
        size: file.size ?? 0
      }));
      console.log('Archivos txt cargados en this.txt:', this.txt);
    },
    error: (err) => {
      console.log('Error al cargar archivos txt:', err);
    }
  });
}
  
  createFoldersFromPath(path: string, id: string): void {
    const parts = path.split('/');
    let currentPath = '';
  
    parts.forEach((part) => {
      // Validar que el segmento no esté vacío o compuesto solo por espacios
      if (part.trim() === '') {
        return; // Ignorar segmentos vacíos
      }
  
      currentPath = currentPath ? `${currentPath}/${part}` : part;
  
      // Verifica si la carpeta ya existe
      if (!this.carpetas.find((folder) => folder.path === currentPath)) {
        this.carpetas.push({ path: currentPath, id, permissions: '664' });
      }
    });
  }

  getCurrentFolders(): { path: string; id: string; permissions: string }[] {
    const depth = this.currentPath ? this.currentPath.split('/').length : 0;
    return this.carpetas.filter((folder) => {
      const parts = folder.path.split('/');
      return (
        parts.length === depth + 1 && // Debe estar en el siguiente nivel
        folder.path.startsWith(this.currentPath) // Debe comenzar con la ruta actual
      );
    });
  }

  openFolder(folder: { path: string; id: string }): void {
    this.currentPath = folder.path; // Actualiza la ruta actual
    this.ruta = this.currentPath.split('/'); // Actualiza la barra de ruta
  }

  goBack(): void {
    if (this.currentPath) {
      const parts = this.currentPath.split('/');
      parts.pop(); // Elimina el último nivel
      this.currentPath = parts.join('/'); // Actualiza la ruta actual
      this.ruta = this.currentPath.split('/'); // Actualiza la barra de ruta
    }
  
    // Si no hay más niveles en la ruta o no se ha ingresado a ninguna carpeta, regresar a las particiones
    if (!this.currentPath) {
      this.volverAParticiones();
    }
  }
  loadDisks(): void {
    this.analyzerService.getDisks().subscribe({
      next: (disks) => {
        this.disks = disks;
      },
      error: (err) => {
        console.error('Error al cargar los discos:', err);
      }
    });
  }
  selectDisk(disk: { name: string }): void {
    this.selectedDisk = disk.name; // Almacena el nombre del disco seleccionado
    this.analyzerService.getPartitions(disk.name).subscribe({
      next: (response) => {
        console.log('Respuesta del backend:', response); // <-- Agrega este log
        if (response.message) {
          this.partitions = []; // Limpia las particiones
          alert(response.message); // Muestra el mensaje de "No existen particiones"
        } else {
          this.partitions = response.partitions; // Carga las particiones desde el backend
        }
      },
      error: (err) => {
        console.error('Error al cargar las particiones:', err);
      }
    });
  }
  selectpartitions(partition: { name: string; id: string }): void {
    this.selectedpart = partition.name; // Almacena el nombre de la partición seleccionada
    this.selectedpartId = partition.id
    console.log('Partición seleccionada:', this.selectedpart);
    this.currentPath = ''; // Reinicia la ruta actual
    this.ruta = []; // Reinicia la barra de ruta
  
    // Cargar carpetas desde el backend
    this.analyzerService.getFolders().subscribe({
      next: (response) => {
        // Filtrar carpetas por el ID de la partición seleccionada
        const filteredFolders = response.filter((folder) => folder.id === partition.id);
  
        // Procesar los paths dinámicamente
        this.carpetas = [];
        filteredFolders.forEach((folder) => {
          this.createFoldersFromPath(folder.path, folder.id);
        });
  
        console.log('Carpetas procesadas:', this.carpetas); // Verifica las carpetas procesadas
      },
      error: (err) => {
        console.error('Error al cargar las carpetas:', err);
      }
    });
    // Cargar archivos txt reales desde el backend
    this.loadTxtFiles();
  }

  volver(): void {
    this.volverEvent.emit(); // Emitir evento para manejar el botón "Volver"
  }
  volverADiscos(): void {
    this.selectedDisk = null; // Limpia el disco seleccionado
    this.selectedpart = null; // Limpia la partición seleccionada
    this.partitions = []; // Limpia las particiones
    this.carpetas = []; // Limpia las carpetas
  }
  volverAParticiones(): void {
    this.selectedpart = null; // Limpia la partición seleccionada
    this.carpetas = []; // Limpia las carpetas
  }
  showNotaModal = false;
  //para los archivos.txt
  selectFile(file: { path: string; contenido: string }): void {
    this.selectedFileContent = file.contenido;
    this.showNotaModal = true;
  }
  closeNotaModal(): void {
    this.showNotaModal = false;
  }

  createFileFromPath(path: string, id: string, contenido: string, size: number): void {
    // Buscar si ya existe el archivo con ese path e id
    const existingFile = this.txt.find(file => file.path === path && file.id === id);
    if (existingFile) {
      // Si existe, solo actualiza el contenido
      existingFile.contenido = contenido;
      console.log(`Archivo ${path} reescrito en partición ${id}`);
      return;
    }
    // Si no existe, lo crea normalmente
    const parts = path.split('/');
    const fileName = parts.pop();
    let currentPath = '';
    parts.forEach((part) => {
      if (part.trim() === '') return;
      currentPath = currentPath ? `${currentPath}/${part}` : part;
      if (!this.carpetas.find((folder) => folder.path === currentPath)) {
        this.carpetas.push({ path: currentPath, id, permissions: '664' });
      }
    });
    if (fileName) {
      this.txt.push({
        path: currentPath ? `${currentPath}/${fileName}` : `/${fileName}`,
        id,
        permissions: '664',
        contenido,
        size,
      });
      console.log(`Archivo creado: ${fileName} en ${currentPath || 'raíz'}`);
    } else {
      console.error('El nombre del archivo no es válido.');
    }
  }

  get filteredTxt(): { path: string; id: string; permissions: string; contenido: string; size: number }[] {
    // Normaliza la ruta actual quitando la barra inicial y final
    const normalizedCurrentPath = this.currentPath.replace(/^\/|\/$/g, '');
    const prefix = normalizedCurrentPath ? normalizedCurrentPath + '/' : '';
    const depth = normalizedCurrentPath
      ? normalizedCurrentPath.split('/').length + 1
      : 1;
    return this.txt.filter((file) => {
      if (file.id !== this.selectedpartId) return false;
      // Normaliza el path del archivo quitando la barra inicial
      const normalizedFilePath = file.path.replace(/^\//, '');
      const fileParts = normalizedFilePath.split('/');
      if (!normalizedCurrentPath) {
        // Raíz: solo archivos con un segmento
        return fileParts.length === 1;
      }
      // Subcarpeta: debe empezar con el prefijo y tener la profundidad correcta
      return normalizedFilePath.startsWith(prefix) && fileParts.length === depth;
    });
  }

  
      

  
}
