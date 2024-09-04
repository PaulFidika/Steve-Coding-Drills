import React, { useEffect } from 'react';
import reactflow, { Handle, Position } from 'reactflow';
const LazyComponent = () => {
    // React.useEffect(() => {
    //     console.log('LazyComponent1 was rendered');
    // }, []);
    return (React.createElement("div", { style: { padding: '10px', border: '1px solid #ddd', borderRadius: '5px', background: '#f0f0f0' } },
        React.createElement(reactflow.Handle, { type: "target", position: reactflow.Position.Top }),
        React.createElement("h2", null, "This is a lazily loaded React Flow node!"),
        React.createElement(reactflow.Handle, { type: "source", position: reactflow.Position.Bottom })));
};
export default LazyComponent;
