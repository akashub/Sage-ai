# Sage AI Frontend

## Sprint 4 Updates

### New Features
1. **Enhanced Navigation**
   - Modernized navigation bar with improved styling
   - Added Docs and About pages with proper routing
   - Clickable logo that redirects to home page

2. **Documentation Page**
   - Comprehensive documentation layout
   - Quick links sidebar for easy navigation
   - Detailed sections for features and API reference
   - Support section with contact information

3. **About Page**
   - Professional team showcase
   - Feature highlights with icons
   - Mission statement section
   - Call-to-action section with gradient background

4. **Chat Interface Improvements**
   - Enhanced ChatSidebar with profile button
   - Smooth transitions and animations
   - Improved message display
   - Better mobile responsiveness

### Technical Updates
1. **Component Structure**
   - Reorganized component hierarchy
   - Improved state management
   - Enhanced error handling
   - Better prop typing

2. **Testing**
   - Added comprehensive Jest tests
   - Component testing for About page
   - Navigation testing
   - UI interaction testing

3. **Styling**
   - Consistent Material-UI theming
   - Improved color scheme
   - Better spacing and layout
   - Enhanced mobile responsiveness

### Getting Started

1. **Prerequisites**
   ```bash
   Node.js >= 14.x
   npm >= 6.x
   ```

2. **Installation**
   ```bash
   # Clone the repository
   git clone https://github.com/your-org/sage-ai.git
   cd sage-ai/v2/frontend

   # Install dependencies
   npm install
   ```

3. **Development**
   ```bash
   # Start development server
   npm start

   # Run tests
   npm test

   # Build for production
   npm run build
   ```

### Project Structure
```
v2/frontend/
├── src/
│   ├── components/
│   │   ├── chat/
│   │   │   ├── ChatSidebar.jsx
│   │   │   ├── ChatWindow.jsx
│   │   │   └── Message.jsx
│   │   ├── layout/
│   │   │   ├── Navigation.jsx
│   │   │   └── Footer.jsx
│   │   └── common/
│   │       ├── Button.jsx
│   │       └── Card.jsx
│   ├── pages/
│   │   ├── Home.jsx
│   │   ├── Docs.jsx
│   │   └── About.jsx
│   ├── __tests__/
│   │   ├── About.test.jsx
│   │   └── Navigation.test.jsx
│   ├── App.jsx
│   └── index.jsx
├── public/
└── package.json
```

### Dependencies
```json
{
  "dependencies": {
    "@emotion/react": "^11.11.0",
    "@emotion/styled": "^11.11.0",
    "@mui/icons-material": "^5.11.16",
    "@mui/material": "^5.13.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.11.1"
  },
  "devDependencies": {
    "@testing-library/jest-dom": "^5.16.5",
    "@testing-library/react": "^14.0.0",
    "jest": "^29.5.0"
  }
}
```

### Testing
Run the test suite:
```bash
npm test
```

Key test files:
- `About.test.jsx`: Tests for About page component
- `Navigation.test.jsx`: Tests for navigation functionality

### Deployment
1. Build the project:
   ```bash
   npm run build
   ```

2. Deploy the contents of the `build` directory to your hosting service.

### Contributing
1. Create a new branch for your feature
2. Make your changes
3. Write tests for new functionality
4. Submit a pull request

### Known Issues
- Mobile responsiveness needs further optimization
- Some edge cases in chat functionality need handling
- Performance optimization needed for large datasets

### Future Improvements
1. **Planned Features**
   - User authentication system
   - Advanced data visualization
   - Custom theme support
   - Export functionality

2. **Technical Debt**
   - Code splitting for better performance
   - Improved error boundaries
   - Enhanced accessibility features
   - Better state management

### Support
For support, email support@sageai.com or create an issue in the repository.

### License
This project is licensed under the MIT License - see the LICENSE file for details. 