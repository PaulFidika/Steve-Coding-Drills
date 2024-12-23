import Surreal from 'surrealdb';
import dotenv from 'dotenv';

dotenv.config();

const requiredEnvVars = [
   'SURREAL_URL',
   'SURREAL_NS', 
   'SURREAL_DB',
   'SURREAL_USER',
   'SURREAL_PASS'
] as const;

// Check if any required env vars are missing
const missingVars = requiredEnvVars.filter(varName => !process.env[varName]);
if (missingVars.length > 0) {
   throw new Error(`Missing required environment variables: ${missingVars.join(', ')}`);
}

export const getDBConnection = async (): Promise<Surreal> => {
   const db = new Surreal();
   
   await db.connect(process.env.SURREAL_URL!, {
       namespace: process.env.SURREAL_NS!,
       database: process.env.SURREAL_DB!,
       auth: {
           username: process.env.SURREAL_USER!,
           password: process.env.SURREAL_PASS!,
       }
   });
   
   return db;
};
