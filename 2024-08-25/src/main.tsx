import React, { Suspense, lazy, useState, useCallback, useEffect } from 'react';
import ReactDOM from 'react-dom/client';
import ReactFlow, { Background, Controls, MiniMap, useReactFlow, Node, ReactFlowProvider } from 'reactflow';
import 'reactflow/dist/style.css';
import PlaceholderNode from './Placeholder';
import { loadModule } from './LoadModule';

const initialNodes = [
  { id: '1', type: "input", position: { x: 0, y: 0 }, data: { label: 'Node 1' } },
  { id: '2', type: "group", position: { x: 0, y: 100 }, data: { label: 'Node 2' } },
];

const initialEdges = [{ id: 'e1-2', source: '1', target: '2' }];

const createLazyNodeType = (url: string) => {
  const LazyComponent = lazy(async () => {
    const module = await loadModule(url);
    return { default: module.default };
  });

  return (props: any) => (
    <Suspense fallback={<PlaceholderNode />}>
      <LazyComponent {...props} />
    </Suspense>
  );
};

const nodeTypes = {
  'lazy-1': createLazyNodeType('http://localhost:3001/Lazy1.js'),
  'lazy-2': createLazyNodeType('http://localhost:3001/Lazy2.js'),
  'lazy-3': createLazyNodeType('http://localhost:3001/Lazy3.js'),
};

const FlowComponent: React.FC = () => {
  const [nodes, setNodes] = useState(initialNodes);
  const { project } = useReactFlow();

  const addNode = useCallback((type: string) => {
    const newNode: Node = {
      id: `${type}-${nodes.length + 1}`,
      type: type,
      position: project({ x: Math.random() * 500, y: Math.random() * 500 }),
      data: { label: `${type} Node` },
    };
    setNodes((nds) => nds.concat(newNode));
  }, [nodes, project]);

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <div style={{ padding: '10px', display: 'flex', gap: '10px' }}>
        <button onClick={() => addNode('lazy-1')}>Add Lazy Node 1</button>
        <button onClick={() => addNode('lazy-2')}>Add Lazy Node 2</button>
        <button onClick={() => addNode('lazy-3')}>Add Lazy Node 3</button>
      </div>

        <ReactFlow 
          nodes={nodes}
          defaultEdges={initialEdges}
          nodeTypes={nodeTypes}
        >
          <Background />
          <Controls />
          <MiniMap />
        </ReactFlow>

    </div>
  );
};

const App: React.FC = () => {
  return (
    <ReactFlowProvider>
      <FlowComponent />
    </ReactFlowProvider>
  );
};

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
