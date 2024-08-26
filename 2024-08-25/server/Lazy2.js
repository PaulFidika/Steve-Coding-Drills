import React, { useState } from 'react';
import { Handle, Position } from 'reactflow';
const LazyComponent = () => {
    const [isExpanded, setIsExpanded] = useState(false);
    console.log('LazyComponent2 was rendered');
    return (React.createElement("div", { style: { padding: '15px', border: '2px solid #007bff', borderRadius: '10px', background: 'linear-gradient(to bottom right, #e6f2ff, #ffffff)' } },
        React.createElement(Handle, { type: "target", position: Position.Top, style: { background: '#007bff' } }),
        React.createElement("h2", { style: { color: '#007bff', textAlign: 'center' } }, "Interactive React Flow Node"),
        React.createElement("p", { style: { textAlign: 'center' } }, "Click to toggle content"),
        isExpanded && (React.createElement("ul", { style: { listStyleType: 'none', padding: 0 } },
            React.createElement("li", null, "Feature 1"),
            React.createElement("li", null, "Feature 2"),
            React.createElement("li", null, "Feature 3"))),
        React.createElement("button", { onClick: () => setIsExpanded(!isExpanded), style: {
                display: 'block',
                margin: '10px auto',
                padding: '5px 10px',
                background: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '5px',
                cursor: 'pointer'
            } }, isExpanded ? 'Collapse' : 'Expand'),
        React.createElement(Handle, { type: "source", position: Position.Bottom, style: { background: '#007bff' } })));
};
export default LazyComponent;
