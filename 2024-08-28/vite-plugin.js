import { exec } from 'child_process';

export default function customBuildPlugin() {
    return {
        name: 'custom-build-plugin',
        buildStart() {
            console.log('Build started...');
        },
        buildEnd() {
            console.log('Build ended...');
        },
        closeBundle() {
            console.log('Bundle closed...');
            // Execute an arbitrary command
            exec('echo "Custom command executed!"', (error, stdout, stderr) => {
                if (error) {
                    console.error(`Error executing command: ${error.message}`);
                    return;
                }
                if (stderr) {
                    console.error(`stderr: ${stderr}`);
                    return;
                }
                console.log(`stdout: ${stdout}`);
            });
        }
    };
}