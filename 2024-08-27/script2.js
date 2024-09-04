import { Whatever } from './script.js'; // circular import
import { SomethingElse } from './script3.js';

// Function to create and append elements
export function createAndAppendElement(tag, text, parent) {
    const element = document.createElement(tag);
    element.textContent = text;
    parent.appendChild(element);
}

// Add message at the bottom of the document
document.addEventListener('DOMContentLoaded', function() {
    const message = document.createElement('p');
    message.textContent = 'script 2 was here';
    document.body.appendChild(message);
});