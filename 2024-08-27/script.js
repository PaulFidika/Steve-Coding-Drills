import { createAndAppendElement } from './script2.js';
// import { Body } from 'easy';

// Define and export Whatever
export const Whatever = {
    name: 'Whatever',
    description: 'This is an exported object called Whatever',
    doSomething: function() {
        console.log('Whatever is doing something!');
    }
};

// Wait for the DOM to be fully loaded
document.addEventListener('DOMContentLoaded', function() {
    const body = document.body;

    // Create a header
    const header = document.createElement('header');
    createAndAppendElement('h1', 'Welcome to My Website', header);
    body.appendChild(header);

    // Create a main content section
    const main = document.createElement('main');
    createAndAppendElement('p', 'This is some basic content rendered using JavaScript.', main);
    
    // Create a list
    const ul = document.createElement('ul');
    ['Item 1', 'Item 2', 'Item 3'].forEach(item => {
        createAndAppendElement('li', item, ul);
    });
    main.appendChild(ul);

    body.appendChild(main);

    // Create a footer
    const footer = document.createElement('footer');
    createAndAppendElement('p', '© 2024 My Website. All rights reserved.', footer);
    body.appendChild(footer);
});
