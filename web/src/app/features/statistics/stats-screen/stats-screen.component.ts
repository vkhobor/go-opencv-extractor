import { Component, computed } from '@angular/core';
import { StatsService } from '../../../services/stats.service';
import { JsonPipe } from '@angular/common';
import { LayoutComponent } from '../../../components/layout/layout.component';

@Component({
  selector: 'app-stats-screen',
  standalone: true,
  imports: [JsonPipe, LayoutComponent],
  templateUrl: './stats-screen.component.html',
  styleUrl: './stats-screen.component.css',
})
export class StatsScreenComponent {
  constructor(private stats: StatsService) {}
  data = this.stats.getStats().result;
}
