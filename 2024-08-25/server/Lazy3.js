import React from 'react';
import { Handle, Position } from 'reactflow';
const LazyComponent = () => {
    console.log('LazyComponent3 was rendered');
    return (React.createElement("div", { style: {
            padding: '15px',
            border: '2px dashed #ff6b6b',
            borderRadius: '15px',
            background: 'linear-gradient(45deg, #ffeaa7, #fab1a0)',
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
            maxWidth: '250px'
        } },
        React.createElement(Handle, { type: "target", position: Position.Top, style: { background: '#ff6b6b' } }),
        React.createElement("h2", { style: { color: '#2d3436', textAlign: 'center', marginBottom: '10px' } }, "Interactive Node"),
        React.createElement("p", { style: { textAlign: 'center', color: '#636e72' } }, "Hover to reveal info"),
        React.createElement("div", { style: {
                overflow: 'hidden',
                maxHeight: '0',
                transition: 'max-height 0.3s ease-out'
            } },
            React.createElement("ul", { style: { listStyleType: 'none', padding: 0, margin: '10px 0' } },
                React.createElement("li", null, "\uD83D\uDE80 Feature A"),
                React.createElement("li", null, "\uD83D\uDCA1 Feature B"),
                React.createElement("li", null, "\uD83D\uDD27 Feature C"))),
        React.createElement(Handle, { type: "source", position: Position.Bottom, style: { background: '#ff6b6b' } })));
};
export default LazyComponent;
