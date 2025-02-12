# 🚀 Sprint 1: Core Implementation of Sage-AI v2

## 🎯 Overview
Sage-AI v2 represents a significant evolution in natural language to SQL query conversion. By combining the robustness of Go with the AI capabilities of Python, we've created a system that not only translates natural language queries but also learns and improves from each interaction.

## 🏗️ Architecture and Components

### High-Level Architecture
![High Level Architecture Diagram](./Sage-ai-v2-Highleveldiagram.png)

Our system follows a hybrid architecture where the core business logic resides in Go while leveraging Python's machine learning capabilities through a bridge pattern. The system comprises three main layers:

1. 🛠️ Core Processing Layer (Go)
   - Handles request processing
   - Manages business logic
   - Coordinates between components
   - Ensures data consistency

2. 🧠 LLM Integration Layer (Python)
   - Processes natural language
   - Generates SQL queries
   - Validates query structure
   - Provides healing capabilities

3. 🎭 Orchestration Layer (Go)
   - Manages processing flow
   - Handles state transitions
   - Coordinates healing attempts
   - Ensures data integrity

### Low-Level Component Interaction
![Low Level Flow Diagram](./Sage-ai-v2-lowleveldiagram.png)

The system processes queries through a sophisticated pipeline:

1. 🔍 Query Analysis
   - Natural language understanding
   - Schema context integration
   - Intent identification
   - Column mapping

2. ⚙️ Query Generation
   - SQL structure creation
   - Schema validation
   - Type checking
   - Syntax verification

3. ✅ Query Validation
   - Structural validation
   - Semantic checking
   - Performance analysis
   - Security verification

4. 🎬 Query Execution
   - Data retrieval
   - Result formatting
   - Error handling
   - Response generation

## 📚 User Stories

1. 💬 Natural Language Query Processing
   "As a data analyst, I want to query my CSV data using natural language so that I don't need to write complex SQL queries manually."
   - Implementation: Created an intuitive CLI interface
   - Achievement: Successfully processes natural language queries
   - Outcome: Generates accurate SQL queries

2. 📊 Schema Understanding
   "As a user, I want the system to understand my CSV structure automatically so that I don't need to provide schema information manually."
   - Implementation: Automatic schema detection
   - Achievement: Correctly identifies column types
   - Outcome: Provides context-aware query generation

3. 🔧 Query Validation and Healing
   "As a developer, I want the system to validate and fix queries automatically so that I get reliable results."
   - Implementation: Multi-stage validation
   - Achievement: Identifies and corrects query issues
   - Outcome: Ensures query reliability

4. 📂 Multi-Dataset Support
   "As a data scientist, I want to use the system with different datasets so that I have a unified query interface."
   - Implementation: Dynamic schema handling
   - Achievement: Works with various CSV formats
   - Outcome: Provides dataset flexibility

5. 🎯 Intelligent Query Generation
   "As an analyst, I want the system to understand complex query requirements so that I can get accurate results."
   - Implementation: Context-aware query generation
   - Achievement: Handles complex query patterns
   - Outcome: Produces precise SQL queries

## 🎯 What Issues Did Our Team Plan to Address

Our team set out to create a system inspired by Vanna.ai but with several ambitious enhancements:

1. 🧠 Knowledge Base Integration
   - Planned to implement a vector database using Qdrant
   - Aimed to store successful query pairs
   - Intended to use historical queries as context
   - Designed for continuous learning

2. 🎭 Advanced Orchestration
   - Developed a sophisticated state management system
   - Implemented healing capabilities
   - Created context-aware processing
   - Designed for extensibility

## ✅ What We Successfully Accomplished

1. 🏗️ Core Architecture Implementation
   - Successfully created Go-Python bridge
   - Implemented robust error handling
   - Developed state management
   - Created logging system

2. ⚙️ Query Processing Pipeline
   - Implemented analysis node
   - Created generation capabilities
   - Added validation system
   - Developed execution handling

3. 💻 CLI Interface
   - Created user-friendly interface
   - Implemented session management
   - Added debug information
   - Developed error reporting

## 🚧 Challenges and Unmet Goals

The transition from Python to Go presented significant challenges:

1. 🔧 Technical Challenges
   - Learning Go's concurrency model
   - Understanding Go's type system
   - Implementing proper error handling
   - Managing memory efficiently

2. 🔄 Integration Complexities
   - Creating effective bridge pattern
   - Handling cross-language communication
   - Managing state across services
   - Ensuring type safety

3. 📚 Knowledge Base Implementation
   - Vector database integration postponed
   - Historical query learning delayed
   - Context enhancement pending
   - Continuous learning features deferred

The shift from v1's pure Python implementation to v2's Go-based architecture, while technically superior, introduced significant complexity. Following Professor Dobra's guidance, we prioritized robustness and maintainability over immediate feature completeness. This decision, while challenging due to our team's limited Go experience, has established a more solid foundation for future enhancements.

Our team had to invest considerable time in learning Go's unique features and best practices, which impacted our ability to implement all planned features. However, the resulting architecture provides better performance, improved error handling, and more robust state management than our original Python implementation.

## 📝 Contributors
- Backend + LLM: Aakash Singh
- Backend + LLM: Nitin Reddy
- Frontend: Yash Kishore
- Frontend: Sudiksha Rajavaram

## 📋 License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.