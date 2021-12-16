const fs = require('fs');
const readline = require('readline');

const { performance, PerformanceObserver } = require('perf_hooks');

const FILENAME = 'realhuman_phill.txt';

async function process() {
  const perfObserver = new PerformanceObserver(items => {
    items.getEntries().forEach(entry => {
      console.log('\n');
      console.log('Duration (in seconds): %f sec\n', (entry.duration / 1000).toFixed(2));
    })
  });

  perfObserver.observe({ entryTypes: ['measure'], buffer: true });

  const stream = fs.createReadStream(FILENAME);

  let lineCount = 0;
  let occurences = 0;

  const result = readline.createInterface({
    input: stream,
    crlfDelay: Infinity,
  });

  performance.mark('read-start');

  for await (const line of result) {
    if (line.includes('japierdole')) {
      console.log(line);

      occurences += 1;
    }

    lineCount += 1;
  }

  performance.mark('read-end');

  performance.measure('read', 'read-start', 'read-end');

  console.log(`Lines count: ${lineCount}`);
  console.log(`Occurences: ${occurences}`);
}

process();
