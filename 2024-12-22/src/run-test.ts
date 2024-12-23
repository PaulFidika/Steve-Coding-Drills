import { Surreal } from 'surrealdb';
import { getDBConnection } from './config';
import { performance } from 'perf_hooks';

interface TestResults {
   operation: string;
   success: boolean;
   duration: number;
}

function pickRandomId(): number {
   return Math.floor(Math.random() * 50000);
}

async function runOperation(db: Surreal, operation: string): Promise<TestResults> {
   const start = performance.now();

   try {
       // Create a promise that rejects after 300ms
       const timeoutPromise = new Promise((_, reject) => {
           setTimeout(() => reject(new Error('Operation timed out')), 10000);
       });

       // Run the operation with timeout
       await Promise.race([
           (async () => {
               switch (operation) {
                   case 'create':
                        const id = pickRandomId();
                        await db.query(`
                            CREATE users:${id} SET
                                name = 'New User ${id}',
                                email = 'user${id}@example.com',
                                power_level = ${Math.floor(Math.random() * 9000) + 1},
                                created_at = time::now()
                        `);
                        break;
                       
                    case 'read':
                        // For some reason this only uses an index ON ASC
                        // and not ORDER BY DESC.
                        await db.query(`
                            SELECT *
                            FROM users
                            ORDER BY power_level ASC
                            LIMIT 20
                        `);
                        break;
                       
                    case 'update':
                        await db.query(`
                            UPDATE users:${pickRandomId()} SET
                                name = 'Updated User ${Math.random()}'
                        `);
                        break;
                       
                    case 'delete':
                        await db.query(`DELETE users:${pickRandomId()}`);
                        break;
                    }
           })(),
           timeoutPromise
       ]);
       
       return {
           operation,
           success: true,
           duration: performance.now() - start
       };
   } catch (error) {
       return {
           operation,
           success: false,
           duration: performance.now() - start
       };
   }
}

async function loadTest(concurrency: number, duration: number) {
   const db = await getDBConnection();
   const operations = ['create', 'read', 'update', 'delete'];
   const results: TestResults[] = [];
   const endTime = Date.now() + (duration * 1000);
   
   async function worker() {
       while (Date.now() < endTime) {
           const operation = operations[Math.floor(Math.random() * operations.length)];
           const result = await runOperation(db, operation);
           results.push(result);
       }
   }
   
   const workers = Array(concurrency).fill(null).map(() => worker());
   await Promise.all(workers);
   
   // Calculate statistics
   const stats = operations.reduce((acc, op) => {
       const opResults = results.filter(r => r.operation === op);
       const successful = opResults.filter(r => r.success);
       const avgDuration = successful.reduce((sum, r) => sum + r.duration, 0) / successful.length;
       
       acc[op] = {
           total: opResults.length,
           successful: successful.length,
           avgDuration: avgDuration.toFixed(2),
           successRate: ((successful.length / opResults.length) * 100).toFixed(2) + '%'
       };
       return acc;
   }, {} as any);
   
   console.log('\nLoad Test Results:');
   console.log('==================');
   console.log(`Concurrency: ${concurrency}`);
   console.log(`Duration: ${duration} seconds`);
   console.log('\nOperation Statistics:');
   console.table(stats);
   
   await db.close();
}

// Run the load test with command line arguments
const concurrency = parseInt(process.argv[2]) || 50;
const duration = parseInt(process.argv[3]) || 30;
loadTest(concurrency, duration);
