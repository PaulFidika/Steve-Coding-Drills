import React, { useEffect } from 'react';
import { Handle, Position } from 'reactflow';
const LazyComponent = () => {
    useEffect(() => {
        console.log('LazyComponent1 was rendered');
    }, []);
    return (React.createElement("div", { style: { padding: '10px', border: '1px solid #ddd', borderRadius: '5px', background: '#f0f0f0' } },
        React.createElement(Handle, { type: "target", position: Position.Top }),
        React.createElement("h2", null, "This is a lazily loaded React Flow node!"),
        React.createElement(Handle, { type: "source", position: Position.Bottom })));
};
export default LazyComponent;
