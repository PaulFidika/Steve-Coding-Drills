import React from 'react';
import { Handle, Position } from 'reactflow';

const LazyComponent: React.FC = () => {
  console.log('LazyComponent3 was rendered');

  return (
    <div style={{ 
      padding: '15px', 
      border: '2px dashed #ff6b6b', 
      borderRadius: '15px', 
      background: 'linear-gradient(45deg, #ffeaa7, #fab1a0)',
      boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
      maxWidth: '250px'
    }}>
      <Handle type="target" position={Position.Top} style={{ background: '#ff6b6b' }} />
      <h2 style={{ color: '#2d3436', textAlign: 'center', marginBottom: '10px' }}>Interactive Node</h2>
      <p style={{ textAlign: 'center', color: '#636e72' }}>Hover to reveal info</p>
      <div style={{ 
        overflow: 'hidden', 
        maxHeight: '0',
        transition: 'max-height 0.3s ease-out'
      }}>
        <ul style={{ listStyleType: 'none', padding: 0, margin: '10px 0' }}>
          <li>🚀 Feature A</li>
          <li>💡 Feature B</li>
          <li>🔧 Feature C</li>
        </ul>
      </div>
      <Handle type="source" position={Position.Bottom} style={{ background: '#ff6b6b' }} />
    </div>
  );
};

export default LazyComponent;
