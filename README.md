# Assessment Management System Documentation

Welcome to the Assessment Management System documentation. This document serves as an entry point to all aspects of the system's architecture, API endpoints, database schema, and more.

## Overview

The Assessment Management System is a comprehensive platform designed to manage organizations, users, courses, and assessments. The system implements Role-Based Access Control (RBAC) with three roles:

1. **Admin** - Can manage all aspects of the system including organizations, users, courses, and assessments
2. **Teacher** - Can manage courses they are assigned to, create assessments, and grade student submissions
3. **Student** - Can enroll in courses, view and submit assessments, and view their grades

## System Architecture

The system follows a clean architecture pattern with the following layers:

1. **Handlers** - HTTP request handlers that parse requests and call appropriate services
2. **Services** - Business logic layer that implements core functionality
3. **Repositories** - Data access layer that interacts with the database
4. **Models** - Data structures used throughout the application

## Key Features

- **Multi-organization Support** - System supports multiple organizations with isolated data
- **Role-Based Access Control** - Granular permissions based on user roles
- **Course Management** - Create, update, and manage courses
- **Assessment System** - Create assessments, submit answers, and grade submissions
- **User Management** - Manage users with different roles
- **Email Uniqueness** - Ensures global email uniqueness across all organizations

- More Details on "Documentations" Folder.
