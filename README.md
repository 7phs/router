# Destination Measurement

A service is responsible for collecting information from 3rd party map-services and decorating it into internal domain model representation.

**Functional Requirements:**

1. Handles REST API requests with a list of geographic coordinates and orders these points by time to travel to them.
2. The service should support easy switching to new services.

**Non-functional Requirements:**

1. Throttles requests to 3rd party map-services to align with their limits.
2. Response time should be less than 20 ms on average for a working day in a specific location.
3. Needs to respond with already known destinations even if 3rd party services are not available at the current moment.

# Architecture

```plantuml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

System(client, "A client of this service", "Services, web- or mobile+ applications", $sprite="clients")

Container_Boundary(api, "API Application") {
    Component(rest_api, "REST API server")
    Component(bridge, "A bridge to map services")
    Component(cache, "A cache of data from map services")
    Component(decorator, "A decorator of map services")
        
    Rel(rest_api, bridge, "Fetch data")
    Rel(bridge, cache, "Load/store cached data")
    Rel(bridge, decorator, "Fetch destination data")
}

System_Ext(osrm, "Mainframe Banking System", "Stores all of the core banking information about customers, accounts, transactions, etc.")

Rel(decorator, osrm, "Fetch geo-data")
@enduml
```