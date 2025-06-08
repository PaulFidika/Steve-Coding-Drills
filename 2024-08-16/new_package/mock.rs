use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::{Duration, Instant};

fn main() {
    let db_data = vec!["id1", "id2", "id3", "id4", "id5", "id6"];
    let mut db_data_map = HashMap::new();
    db_data_map.insert("id1", "RandomString1");
    db_data_map.insert("id2", "AnotherRandomString");
    db_data_map.insert("id3", "YetAnotherString");
    db_data_map.insert("id4", "FourthRandomString");
    db_data_map.insert("id5", "FifthRandomString");
    db_data_map.insert("id6", "LastRandomString");

    let db_data_map = Arc::new(db_data_map);
    let expected_results = db_data.len();
    let my_counter = Arc::new(Mutex::new(1));

    let (tx, rx) = std::sync::mpsc::channel();

    let t0 = Instant::now();

    let handles: Vec<_> = db_data
        .into_iter()
        .map(|id| {
            let tx = tx.clone();
            let db_data_map = Arc::clone(&db_data_map);
            let my_counter = Arc::clone(&my_counter);

            thread::spawn(move || {
                tx.send(db_call(&id, &db_data_map)).unwrap();
                
                // Without this lock, this block has inconsistent behavior
                // let mut counter = my_counter.lock().unwrap();
                let thing = *my_counter.lock().unwrap();
                thread::sleep(Duration::from_nanos(100_000_000));
                *my_counter.lock().unwrap() = thing + 1;
                // drop(counter);
            })
        })
        .collect();

    let mut i = 0;
    for result in rx.iter().take(expected_results) {
        println!("Received result: {}", result);
        i += 1;
        if i == expected_results {
            break;
        }
    }

    for handle in handles {
        handle.join().unwrap();
    }

    println!("\nCounter results: {}", *my_counter.lock().unwrap());
    println!("\nTotal execution time: {:?}", t0.elapsed());
}

fn db_call(id: &str, db_data_map: &HashMap<&str, &str>) -> String {
    let delay = Duration::from_millis(2000);
    thread::sleep(delay);
    println!("The result from the database is: {}", id);
    db_data_map[id].to_string()
}