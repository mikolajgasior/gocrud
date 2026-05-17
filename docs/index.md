# Overview

**gocrud** is a Go library designed to streamline the development of CRUDL (Create, Read, Update, Delete, List) APIs 
and interfaces for PostgreSQL-backed applications. By leveraging reflection, gocrud aims to automatically generate 
database persistence and API layers from Go structs containing basic field types (such as int, string, time.Time, etc.).

The library is intended to provide:
* Database persistence layer: Maps Go structs to PostgreSQL tables, handling schema creation and data operations.
* RESTful API: Generates a ready-to-use CRUD API endpoint set for your struct.
* Admin UI: A simple web interface to browse and manage records (currently under development).

## Motivation

**gocrud** was created to address the need for rapidly building APIs during prototyping phases for internal customer 
systems. The goal is to enable developers to go from a struct definition to a functional data management system with 
minimal configuration—ideal for prototypes, internal tools, MVPs, or any scenario where speed and simplicity matter.

## Demo

Please navigate to [Demo](demo.md) see sample usage.