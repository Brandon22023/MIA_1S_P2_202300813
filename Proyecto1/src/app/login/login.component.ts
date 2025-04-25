import { Component, EventEmitter, Output } from '@angular/core';
import { AnalyzerService } from '../analyzer.service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  imports: [FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
})
export class LoginComponent {
  partitionId: string = ''; // ID de la partición
  username: string = ''; // Usuario
  password: string = ''; // Contraseña
  constructor(private analyzerService: AnalyzerService) {}

  iniciarSesion(): void {
    this.analyzerService.login(this.username, this.password, this.partitionId).subscribe({
      next: (response) => {
        console.log('Respuesta del servidor:', response);
        alert(response.message || 'Inicio de sesión exitoso');
      },
      error: (error) => {
        console.error('Error del servidor:', error);
        alert(error.error?.message || 'Error al iniciar sesión');
      },
    });
  }
  @Output() volverEvent = new EventEmitter<void>();

  volver(): void {
    this.volverEvent.emit(); // Emitir evento para manejar el botón "Volver"
  }
}
