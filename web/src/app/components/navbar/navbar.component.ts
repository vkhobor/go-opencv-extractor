import { Component, Input } from '@angular/core';
import { SettingsToggleComponent } from '../settings-toggle/settings-toggle.component';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [SettingsToggleComponent, RouterLink, RouterLinkActive],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css',
})
export class NavbarComponent {
  @Input() sticky = true;
  @Input() id = 'default-id';
}
