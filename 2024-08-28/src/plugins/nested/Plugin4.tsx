import React, { useState } from 'react';
import { PluginProps } from '../../interface';

interface RandomQuoteProps extends PluginProps {
  author: string;
}

const RandomQuote: React.FC<RandomQuoteProps> = ({ author }) => {
  const quotes = [
    "The only way to do great work is to love what you do.",
    "Innovation distinguishes between a leader and a follower.",
    "Stay hungry, stay foolish.",
    "Your time is limited, don't waste it living someone else's life.",
    "The future belongs to those who believe in the beauty of their dreams."
  ];

  const [quote, setQuote] = useState(quotes[0]);

  const generateRandomQuote = () => {
    const randomIndex = Math.floor(Math.random() * quotes.length);
    setQuote(quotes[randomIndex]);
  };

  return (
    <div>
      <h3>Random Quote by {author}</h3>
      <p>{quote}</p>
      <button onClick={generateRandomQuote}>New Quote</button>
    </div>
  );
};

export default RandomQuote;
