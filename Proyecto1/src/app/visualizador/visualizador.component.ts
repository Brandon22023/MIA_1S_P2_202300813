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
  selectedDisk: string | null = null;
  constructor(private analyzerService: AnalyzerService) {}

  ngOnInit(): void {
    this.loadDisks();
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

  volver(): void {
    this.volverEvent.emit(); // Emitir evento para manejar el botón "Volver"
  }
  volverADiscos(): void {
    this.selectedDisk = null; // Regresa a la vista de discos
  }

}
