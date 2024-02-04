# Destination Measurement

A service is responsible for collecting information from 3rd party map-services and decorating it into internal domain model representation.

**Functional Requirements:**

1. Handles REST API requests with a list of geographic coordinates and orders these points by time to travel to them.
2. The service should support easy switsching to new services.

**Non-functional Requirements:**

1. Throttles requests to 3rd party map-services to align with their limits.
2. Response time should be less than 20 ms on average for a working day in a specific location.
3. Needs to respond with already known destinations even if 3rd party services are not available at the current moment.

# Architecture

```plantuml
@startuml architecture
!if %variable_exists("RELATIVE_INCLUDE")
    !include %get_variable_value("RELATIVE_INCLUDE")/C4_Context.puml
    !include %get_variable_value("RELATIVE_INCLUDE")/C4_Container.puml
    !include %get_variable_value("RELATIVE_INCLUDE")/C4_Component.puml
!else
    !include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml
    !include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml
    !include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml
!endif

SHOW_PERSON_OUTLINE()

AddElementTag("proposed", $bgColor="#666666", $fontColor="#bfbfbf", $borderColor="#bfbfbf")
AddRelTag("async", $textColor=$ARROW_FONT_COLOR, $lineColor=$ARROW_COLOR, $lineStyle=DashedLine())
AddRelTag("sync/async", $textColor=$ARROW_FONT_COLOR, $lineColor=$ARROW_COLOR, $lineStyle=DottedLine())

title Component diagram for Routing service

System(client, "A client of service", "Services, web- or mobile- applications", $sprite="clients")

System_Boundary(api, "Routing Service") {
    Component(rest_api, "REST API server")
    Component(bridge, "A bridge to routing services")
    ComponentDb(mem_cache, "An in-memory cache")
    Component(cache, "A cache of data from routing services")
    Component(decorator, "A decorator of external routing services")
        
    Rel_Neighbor(cache, mem_cache, "Fetch/store", "sync")
    Rel(rest_api, bridge, "Fetch routing data", "sync")
    Rel(bridge, cache, "Fetch/store cached data", "sync")
    Rel(bridge, decorator, "Fetch rouging data", "async")
}

ContainerDb(cach_db, "Cache Database", "NoSQL", "Stores  with TTL routing data already fetched", "proposed", $sprite="postgresql,color=gray")

System_Ext(osrm, "Routing Machine service")

Rel(client, rest_api, "Fetch geo-data", "sync")
Rel(cache, cach_db, "Fetch cached data", "sync")
Rel(decorator, osrm, "Fetch geo-data", "async")
@enduml
```