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
  carpetas: { path: string; id: string }[] = [];
  currentPath: string = ''; // Ruta actual
  ruta: string[] = []; // Ruta desglosada para mostrar en la barra
  selectedDisk: string | null = null;
  selectedpart: string | null = null;
  constructor(private analyzerService: AnalyzerService) {}

  ngOnInit(): void {
    this.loadDisks();
  }
  createFoldersFromPath(path: string, id: string): void {
    const parts = path.split('/');
    let currentPath = '';

    parts.forEach((part) => {
      currentPath = currentPath ? `${currentPath}/${part}` : part;

      // Verifica si la carpeta ya existe
      if (!this.carpetas.find((folder) => folder.path === currentPath)) {
        this.carpetas.push({ path: currentPath, id });
      }
    });
  }

  getCurrentFolders(): { path: string; id: string }[] {
    // Filtra las carpetas que están en el nivel actual
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
    this.currentPath = ''; // Reinicia la ruta actual
    this.ruta = []; // Reinicia la barra de ruta
  
    // Generar carpetas dinámicamente basadas en el ID de la partición
    if (partition.id === '131A') {
      this.createFoldersFromPath('home/user/pinqui', partition.id);
      this.createFoldersFromPath('home/user/doc', partition.id);
    } else if (partition.id === '132A') {
      this.createFoldersFromPath('home/documents', partition.id);
      this.createFoldersFromPath('home/images', partition.id);
    } else {
      // Si no hay contenido para la partición, limpiar carpetas
      this.carpetas = [];
    }
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
}
