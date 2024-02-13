import { Component, Input } from '@angular/core';
import { SettingsToggleComponent } from '../settings-toggle/settings-toggle.component';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [SettingsToggleComponent],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css',
})
export class NavbarComponent {
  @Input() sticky = true;
  @Input() id = 'default-id';
}
