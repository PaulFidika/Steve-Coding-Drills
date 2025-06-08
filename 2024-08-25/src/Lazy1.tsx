import React, { useEffect, Suspense } from 'react';
import { Handle, Position } from 'reactflow';

const LazyComponent: React.FC = () => {
  useEffect(() => {
    console.log('LazyComponent1 was rendered');
  }, []);

  return (
    <div style={{ padding: '10px', border: '1px solid #ddd', borderRadius: '5px', background: '#f0f0f0' }}>
      <Handle type="target" position={Position.Top} />
      <h2>This is a lazily loaded React Flow node!</h2>
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
};

export default LazyComponent;
