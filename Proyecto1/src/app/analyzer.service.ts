import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AnalyzerService {
  private apiUrl = 'http://localhost:3000/analyze';

  constructor(private http: HttpClient) {}

  analyze(input: string): Observable<any> {
    console.log('Enviando al servidor:', { command: input }); // Verifica el comando enviado
    return this.http.post<any>(this.apiUrl, { command: input });
  }
}
