import React, { useState } from 'react';
import { Handle, Position } from 'reactflow';

const LazyComponent: React.FC = () => {
  const [isExpanded, setIsExpanded] = useState(false);
  console.log('LazyComponent2 was rendered');

  return (
    <div style={{ padding: '15px', border: '2px solid #007bff', borderRadius: '10px', background: 'linear-gradient(to bottom right, #e6f2ff, #ffffff)' }}>
      <Handle type="target" position={Position.Top} style={{ background: '#007bff' }} />
      <h2 style={{ color: '#007bff', textAlign: 'center' }}>Interactive React Flow Node</h2>
      <p style={{ textAlign: 'center' }}>Click to toggle content</p>
      {isExpanded && (
        <ul style={{ listStyleType: 'none', padding: 0 }}>
          <li>Feature 1</li>
          <li>Feature 2</li>
          <li>Feature 3</li>
        </ul>
      )}
      <button 
        onClick={() => setIsExpanded(!isExpanded)}
        style={{ 
          display: 'block', 
          margin: '10px auto', 
          padding: '5px 10px', 
          background: '#007bff', 
          color: 'white', 
          border: 'none', 
          borderRadius: '5px', 
          cursor: 'pointer' 
        }}
      >
        {isExpanded ? 'Collapse' : 'Expand'}
      </button>
      <Handle type="source" position={Position.Bottom} style={{ background: '#007bff' }} />
    </div>
  );
};

export default LazyComponent;