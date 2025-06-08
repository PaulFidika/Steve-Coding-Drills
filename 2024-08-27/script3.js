// Define and export SomethingElse
export const SomethingElse = {
    name: 'SomethingElse',
    description: 'This is an exported object called SomethingElse',
    doSomething: function() {
        console.log('SomethingElse is doing something!');
    }
};

// Add message at the bottom of the document
document.addEventListener('DOMContentLoaded', function() {
    const message = document.createElement('p');
    message.textContent = 'script 3 was here';
    document.body.appendChild(message);
});
