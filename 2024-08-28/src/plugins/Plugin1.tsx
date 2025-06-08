import React from 'react';
import { PluginProps } from '../interface';

interface DemoComponentProps extends PluginProps {
  title: string;
}

const DemoComponent: React.FC<DemoComponentProps> = ({ title }) => {
  return (
    <div>
      <h1>{title}</h1>
      <p>This is a demo React component.</p>
    </div>
  );
};

export default DemoComponent;
