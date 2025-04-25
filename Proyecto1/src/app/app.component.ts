import { Component, ViewChild, ElementRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AnalyzerService } from './analyzer.service';
import { CommonModule } from '@angular/common'; // Importar CommonModule

@Component({
  selector: 'app-root',
  imports: [FormsModule, CommonModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  standalone: true
})
export class AppComponent {
  isModalVisible: boolean = false; // Estado del modal
  isModalClosing: boolean = false; // Estado de cierre del modal
  title = 'Proyecto1';
  entrada: string = '';
  salida: string = '';
  mensaje: string = ''; // esto sera para mostrar el exito de lo metico

  constructor(private analyzerService: AnalyzerService) {}

  @ViewChild('fileInput') fileInput!: ElementRef;

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      const file = input.files[0];
      const reader = new FileReader();
      reader.onload = (e) => {
        const text = reader.result as string;
        this.entrada = text;
      };
      this.entrada = '';
      reader.readAsText(file);
    }
  }

  limpiar(): void {
    this.playClickSound(); // Reproducir sonido al hacer clic en "Limpiar"
    this.entrada = '';
    this.salida = '';
    this.mensaje = ''; // Limpiar el mensaje
  }

  ejecutar(): void {
    this.playClickSound(); // Reproducir sonido al hacer clic en "Ejecutar"
    this.analyzerService.analyze(this.entrada).subscribe({
      next: (response) => {
        // Imprime la respuesta en la consola
        console.log('Respuesta del servidor:', response);
  
        // Asegúrate de que la respuesta tenga la estructura esperada
        if (response && response.output) {
          this.salida = response.output;
          console.log('Salida procesada:', this.salida);
  
          // Verificar si la salida comienza con "Error"
          if (this.salida.trim().toLowerCase().startsWith('error')) {
            // Mostrar el mensaje de error en el modal
            this.mensaje = this.salida;
            this.showModal();
          } else {
            // Mostrar mensaje de éxito si no hay errores
            this.mensaje = 'Todos los comandos se ejecutaron con éxito';
            this.showModal();
          }
        } else {
          // Si la respuesta no tiene la estructura esperada
          this.salida = 'Respuesta inesperada del servidor';
          this.mensaje = this.salida; // Mostrar el mensaje en el modal
          this.showModal();
        }
      },
      error: (error) => {
        // Manejo de errores
        console.error('Error del servidor:', error);
  
        if (error.error && error.error.error) {
          this.salida = `Error: ${error.error.error}`;
        } else if (error.message) {
          this.salida = `Error: ${error.message}`;
        } else {
          this.salida = 'Error desconocido';
        }
  
        // Mostrar el mensaje de error en el modal
        this.mensaje = this.salida;
        this.showModal();
      }
    });
  }

  // Método para reproducir el sonido
  playClickSound(): void {
    const audio = document.getElementById('clickSound') as HTMLAudioElement;
    if (audio) {
      audio.play();
    }
  }

  showModal(): void {
    this.isModalVisible = true; // Mostrar el modal
    this.isModalClosing = false; // Asegurarse de que no esté cerrando
    this.playClickSound(); // Reproducir sonido al abrir el modal
  }

  closeModal(): void {
    this.isModalClosing = true; // Activar animación de salida
    setTimeout(() => {
      this.isModalVisible = false; // Ocultar el modal después de la animación
      this.isModalClosing = false; // Reiniciar el estado de cierre
    }, 300); // Duración de la animación de salida (0.3s)
    this.playClickSound(); // Reproducir sonido al cerrar el modal
  }

}
