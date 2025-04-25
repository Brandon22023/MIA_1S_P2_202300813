import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-visualizador',
  imports: [],
  templateUrl: './visualizador.component.html',
  styleUrl: './visualizador.component.css'
})
export class VisualizadorComponent {

  @Output() volverEvent = new EventEmitter<void>(); // Evento para volver al componente principal

  volver(): void {
    this.volverEvent.emit(); // Emitir evento para manejar el botón "Volver"
  }

}
