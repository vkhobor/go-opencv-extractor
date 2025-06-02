import { Component, OnInit, computed } from "@angular/core";
import { Router, RouterOutlet } from "@angular/router";
import { ButtonModule } from "primeng/button";
import { initFlowbite } from "flowbite";
import { NavbarComponent } from "./components/navbar/navbar.component";
import { JobsService } from "./services/jobs.service";

@Component({
    selector: "app-root",
    imports: [RouterOutlet, NavbarComponent],
    providers: [],
    templateUrl: "./app.component.html",
    styleUrl: "./app.component.css"
})
export class AppComponent implements OnInit {
	ngOnInit(): void {
		// initFlowbite();
	}
}
