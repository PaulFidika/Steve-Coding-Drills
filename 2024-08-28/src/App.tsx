import React, { useEffect, useState } from 'react';

const App: React.FC = () => {
    const [plugins, setPlugins] = useState<React.ComponentType<any>[]>([]);

    useEffect(() => {
        const loadPlugins = async () => {
            const pluginModules = import.meta.glob('./plugins/**/*.tsx');
            const loadedPlugins: React.ComponentType<any>[] = [];

            for (const path in pluginModules) {
                const module = await pluginModules[path]();
                loadedPlugins.push(module.default);
            }

            setPlugins(loadedPlugins);
        };

        loadPlugins();
    }, []);

    return (
        <div>
            <h1>My Vite App with Plugins</h1>
            {plugins.map((Plugin, index) => (
                <div key={index}>
                    <Plugin />
                </div>
            ))}
        </div>
    );
};

export default App;