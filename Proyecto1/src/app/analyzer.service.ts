import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';

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
  login(user: string, pass: string, id: string): Observable<any> {
    const payload = { user, pass, id };
    console.log('Enviando datos de login:', payload); // Verifica los datos enviados
    return this.http.post<any>('http://localhost:3000/login', payload);
  }
  logout(): Observable<any> {
    console.log('Cerrando sesión...');
    return this.http.post<any>('http://localhost:3000/logout', {});
  }
  getDisks(): Observable<{ name: string; size: string; fit: string; mounted_partitions: string }[]> {
    return this.http.get<{ disks: { name: string; size: string; fit: string; mounted_partitions: string }[] }>('http://localhost:3000/disks').pipe(
      map((response) => response.disks)
    );
  }
  getPartitions(diskName: string): Observable<any> {
    return this.http.get<{ partitions?: any[]; message?: string }>(
      `http://localhost:3000/partitions/${diskName}`
    ).pipe(
      map((response) => {
        if (response.message) {
          return { message: response.message }; // Si no hay particiones
        }
        return { partitions: response.partitions }; // Si hay particiones
      })
    );
  }
  getFolders(): Observable<{ path: string; id: string }[]> {
    return this.http.get<{ carpetas: { path: string; id: string }[] }>('http://localhost:3000/folders').pipe(
      map((response) => response.carpetas) // Extrae la propiedad 'carpetas' del objeto de respuesta
    );
  }
  getTxtFiles(): Observable<{ path: string; id: string; contenido: string; size: number }[]> {
    return this.http.get<{ txtfiles: { path: string; id: string; contenido: string; size: number }[] }>('http://localhost:3000/txtfiles')
      .pipe(map(response => response.txtfiles));
  }
}
