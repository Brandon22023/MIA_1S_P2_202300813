import { Component, EventEmitter, Output } from '@angular/core';
import { AnalyzerService } from '../analyzer.service';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-login',
  imports: [FormsModule, CommonModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
})
export class LoginComponent {
  partitionId: string = ''; // ID de la partición
  username: string = ''; // Usuario
  password: string = ''; // Contraseña

  @Output() irVisualizadorEvent = new EventEmitter<void>(); // Evento para cambiar al visualizador
  @Output() volverEvent = new EventEmitter<void>();
  constructor(private analyzerService: AnalyzerService) {}
  
  iniciarSesion(): void {
    this.analyzerService.login(this.username, this.password, this.partitionId).subscribe({
      next: (response) => {
        console.log('Respuesta del servidor:', response);
        alert(response.message || 'Inicio de sesión exitoso');
        this.irVisualizadorEvent.emit(); // Emitir evento para cambiar al visualizador

      },
      error: (error) => {
        console.error('Error del servidor:', error);
        alert(error.error?.message || 'Error al iniciar sesión');
      },
    });
  }
 

  volver(): void {
    this.volverEvent.emit(); // Emitir evento para manejar el botón "Volver"
  }
}
