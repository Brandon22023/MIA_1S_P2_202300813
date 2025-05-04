import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})

export class AnalyzerService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  analyze(input: string): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/analyze`, { command: input });
  }
  login(user: string, pass: string, id: string): Observable<any> {
    const payload = { user, pass, id };
    return this.http.post<any>(`${this.apiUrl}/login`, payload);
  }
  logout(): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/logout`, {});
  }
  getDisks(): Observable<{ name: string; size: string; fit: string; mounted_partitions: string }[]> {
    return this.http.get<{ disks: { name: string; size: string; fit: string; mounted_partitions: string }[] }>(`${this.apiUrl}/disks`).pipe(
      map((response) => response.disks)
    );
  }
  getPartitions(diskName: string): Observable<any> {
    return this.http.get<{ partitions?: any[]; message?: string }>(
      `${this.apiUrl}/partitions/${diskName}`
    ).pipe(
      map((response) => {
        if (response.message) {
          return { message: response.message };
        }
        return { partitions: response.partitions };
      })
    );
  }
  getFolders(): Observable<{ path: string; id: string }[]> {
    return this.http.get<{ carpetas: { path: string; id: string }[] }>(`${this.apiUrl}/folders`).pipe(
      map((response) => response.carpetas)
    );
  }
  getTxtFiles(): Observable<{ path: string; id: string; contenido: string; size: number }[]> {
    return this.http.get<{ txtfiles: { path: string; id: string; contenido: string; size: number }[] }>(`${this.apiUrl}/txtfiles`)
      .pipe(map(response => response.txtfiles));
  }
}
