import { Component, Input } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { Flowbite } from '../../util/flowbiteFix';

@Component({
    selector: 'app-navbar',
    imports: [RouterLink, RouterLinkActive],
    templateUrl: './navbar.component.html',
    styleUrl: './navbar.component.css'
})
@Flowbite()
export class NavbarComponent {
    @Input() sticky = true;
    @Input() id = 'default-id';
}
