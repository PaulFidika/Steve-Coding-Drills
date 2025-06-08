import { getDBConnection } from './config';

async function initializeSchema() {
   const db = await getDBConnection();
   
   try {
       // Create tables and define schemas
      await db.query(`
          DEFINE TABLE users SCHEMAFULL;
          DEFINE FIELD name ON users TYPE string;
          DEFINE FIELD email ON users TYPE string;
          DEFINE FIELD created_at ON users TYPE datetime VALUE $before OR time::now();
          
          DEFINE TABLE posts SCHEMAFULL;
          DEFINE FIELD title ON posts TYPE string;
          DEFINE FIELD content ON posts TYPE string;
          DEFINE FIELD author ON posts TYPE record<users>;
          DEFINE FIELD created_at ON posts TYPE datetime VALUE $before OR time::now();
          
          DEFINE INDEX idx_users_email ON users FIELDS email UNIQUE;
      `);
       
       // Create some initial data
       await db.query(`
           CREATE users:test1 SET 
               name = 'Test User 1',
               email = 'test1@example.com',
               created_at = time::now();
               
           CREATE users:test2 SET 
               name = 'Test User 2',
               email = 'test2@example.com',
               created_at = time::now();
       `);
       
       console.log('Schema and initial data created successfully');
   } catch (error) {
       console.error('Error initializing schema:', error);
   } finally {
       await db.close();
   }
}

initializeSchema();