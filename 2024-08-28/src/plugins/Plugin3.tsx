import React from 'react';
import { PluginProps } from '../interface';

interface Plugin3Props extends PluginProps {
  title: string;
  description: string;
}

const Plugin3: React.FC<Plugin3Props> = ({ title, description }) => {
  return (
    <div>
      <h2>{title}</h2>
      <p>{description}</p>
      <button onClick={() => alert('Plugin 3 clicked!')}>
        Click me
      </button>
    </div>
  );
};

export default Plugin3;
