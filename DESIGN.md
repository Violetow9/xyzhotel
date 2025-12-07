# Design Stratégique - XYZ Hôtel

## 1. Ubiquitous Language

Voici le résumé des termes métier utilisés dans le code et les échanges avec les experts du domaine.

### Acteurs
* **Customer (Client) :** Personne physique possédant un compte et un portefeuille électronique pour réserver des chambres.
* **Admin :** Gestionnaire de l'hôtel ayant accès aux vues d'ensemble (historique, occupation).

### Finances
* **Wallet (Portefeuille) :** Compte prépayé du client.
* **Balance (Solde) :** Montant disponible dans le portefeuille, exprimé en centimes (pour éviter les erreurs d'arrondi).
* **Currency (Devise) :** Monnaie d'entrée (EUR, USD, etc). Le système convertit tout en EUR pour le stockage.
* **Deposit (Acompte) :** Somme correspondant à 50% du prix total, prélevée immédiatement à la réservation.
* **Due Amount (Reste à payer) :** Les 50% restants à payer pour confirmer la réservation.

### Hébergement
* **Room (Chambre) :** Unité locative physique identifiée par un numéro unique (ex: "101").
* **Room Type (Config) :** Catégorie de chambre déterminant le prix et les équipements (Standard, Superior, Suite).
* **Check-In Date :** Date d'arrivée prévue (normalisée à minuit).
* **Nights (Nombre de nuits) :** Durée du séjour.

### Réservation (Cycle de vie)
* **Reservation :** Contrat de location temporaire entre un client et une chambre.
* **PENDING (En attente) :** Réservation créée, acompte de 50% payé.
* **CONFIRMED (Confirmée) :** Solde restant payé. Le client peut accéder à la chambre.
* **CANCELLED (Annulée) :** Réservation annulée par le client (pas de remboursement).
* **COMPLETED (Terminée) :** Le client a quitté la chambre (Checkout effectué).

---

## 2. Bounded Contexts & Context Maps

```mermaid
flowchart TD
    classDef context fill:#f9f9f9,stroke:#616161,stroke-width:2px,stroke-dasharray: 5 5,color:#616161
    classDef command fill:#90caf9,stroke:#1e88e5,stroke-width:1px,color:black,rx:2,ry:2
    classDef entity fill:#fff59d,stroke:#fbc02d,stroke-width:2px,color:black,rx:0,ry:0
    classDef vo fill:#b2dfdb,stroke:#00695c,stroke-width:1px,color:black,rx:5,ry:5
    classDef event fill:#ffcc80,stroke:#f57c00,stroke-width:1px,color:black,rx:0,ry:0
    classDef system fill:#f48fb1,stroke:#d81b60,stroke-width:1px,color:black,rx:0,ry:0
    classDef actor fill:#ffffff,stroke:#000,stroke-width:1px,color:black

    subgraph Legend ["LÉGENDE DU DOMAINE"]
        direction LR
        L_Cmd["Commande"]:::command
        L_Ent["Entité / Agrégat"]:::entity
        L_VO["Value Object"]:::vo
        L_Evt["Événement"]:::event
        L_Sys["Système"]:::system
        L_Act([Acteur]):::actor
    end

    User([Client]):::actor
    Admin([Admin]):::actor

    subgraph CustomerContext ["BOUNDED CONTEXT: CUSTOMER"]
        direction TB
        
        E_Cust("Entité:<br/>Customer"):::entity
        VO_Wall("VO:<br/>Wallet"):::vo
        VO_Prof("VO:<br/>Email/Profile"):::vo
        
        E_Cust -.-> VO_Wall
        E_Cust -.-> VO_Prof

        C_Create["Create<br/>Customer"]:::command --> E_Cust
        E_Cust --> Ev_Created["Account<br/>Created"]:::event

        C_Credit["Credit<br/>Wallet"]:::command --> E_Cust
        E_Cust --> Ev_Credited["Wallet<br/>Credited"]:::event
    end
    class CustomerContext context

    subgraph RoomContext ["BOUNDED CONTEXT: ROOM"]
        direction TB
        
        E_Room("Entité:<br/>Room"):::entity
        VO_Type("VO:<br/>RoomType"):::vo
        VO_Conf("VO:<br/>PriceConfig"):::vo
        
        E_Room -.-> VO_Type
        E_Room -.-> VO_Conf

        C_List["List<br/>Rooms"]:::command --> E_Room
        E_Room --> Ev_List["Rooms<br/>Displayed"]:::event
    end
    class RoomContext context

    subgraph ReservationContext ["BOUNDED CONTEXT: RESERVATION"]
        direction TB
        
        E_Res("Entité:<br/>Reservation"):::entity
        VO_Money("VO:<br/>Money"):::vo
        VO_Date("VO:<br/>Dates"):::vo

        E_Res -.-> VO_Money
        E_Res -.-> VO_Date
        
        C_Book["Book<br/>Room"]:::command --> E_Res
        E_Res --> Ev_Booked["Res Created<br/>(Pending)"]:::event
        Ev_Booked -.-> S_Debit1["Debit 50%<br/>(System)"]:::system

        C_Conf["Confirm<br/>Res"]:::command --> E_Res
        E_Res --> Ev_Conf["Res<br/>Confirmed"]:::event
        Ev_Conf -.-> S_Debit2["Debit Balance<br/>(System)"]:::system

        C_Cancel["Cancel<br/>Res"]:::command --> E_Res
        E_Res --> Ev_Canc["Res<br/>Cancelled"]:::event

        C_Check["Checkout<br/>Room"]:::command --> E_Res
        E_Res --> Ev_Comp["Res<br/>Completed"]:::event

        C_Occ["View<br/>Occupied"]:::command --> S_Occ["Occupancy List<br/>(Read Model)"]:::system
        S_Occ --> Ev_Occ["Occupancy<br/>Displayed"]:::event

        C_Hist["View<br/>History"]:::command --> S_Hist["History List<br/>(Read Model)"]:::system
        S_Hist --> Ev_Hist["History<br/>Displayed"]:::event
    end
    class ReservationContext context

    User --> C_Create
    User --> C_Credit
    User --> C_List
    User --> C_Book
    User --> C_Conf
    User --> C_Cancel
    User --> C_Check

    Admin --> C_Occ
    Admin --> C_Hist

    RoomContext ===>|"Context Map:<br/>Fournit Prix & Infos"| ReservationContext
    CustomerContext ===>|"Context Map:<br/>Fournit Identité & Fonds"| ReservationContext