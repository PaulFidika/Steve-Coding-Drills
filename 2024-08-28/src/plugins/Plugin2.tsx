import React from 'react';
import { PluginProps } from '../interface';

interface Plugin2Props extends PluginProps {
  message: string;
}

const Plugin2: React.FC<Plugin2Props> = ({ message }) => {
  return (
    <div>
      <h2>Plugin 2</h2>
      <p>{message}</p>
    </div>
  );
};

export default Plugin2;
