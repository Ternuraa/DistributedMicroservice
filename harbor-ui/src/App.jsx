import React, { useState } from 'react';
import { listings as initialListings } from './data';
import ListingCard from "./components/ListingCard.jsx";
import BookingModal from "./components/BookingModal.jsx";

export default function App() {
  const [apartments, setApartments] = useState(initialListings);
  const [selectedListing, setSelectedListing] = useState(null);
  const [toast, setToast] = useState({ message: '', colorClass: '', isVisible: false });

  const showToast = (message, colorClass) => {
    setToast({ message, colorClass, isVisible: true });
    setTimeout(() => {
      setToast(prev => ({ ...prev, isVisible: false }));
    }, 3000);
  };

  const handleBookingSuccess = (listingId, start, end) => {
    setApartments(prev => prev.map(apt => {
      if (apt.id === listingId) {
        const newBookedDays = [...apt.bookedDays];
        for (let d = start; d <= end; d++) {
          if (!newBookedDays.includes(d)) newBookedDays.push(d);
        }
        return { ...apt, bookedDays: newBookedDays };
      }
      return apt;
    }));
  };

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-4xl font-bold mb-2 text-center text-indigo-600">Поиск и бронирование</h1>
        <p className="text-gray-500 text-center mb-10">Выберите интересующие апартаменты, чтобы посмотреть свободные даты</p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {apartments.map(listing => (
            <ListingCard 
              key={listing.id} 
              listing={listing} 
              onClick={(item) => setSelectedListing(item)} 
            />
          ))}
        </div>
      </div>

      {selectedListing && (
        <BookingModal 
          listing={selectedListing} 
          onClose={() => setSelectedListing(null)} 
          onSuccess={handleBookingSuccess}
          showToast={showToast}
        />
      )}

      <div 
        className={`fixed bottom-5 right-5 text-white px-6 py-3 rounded-xl shadow-xl transition-all duration-300 transform z-50 ${toast.colorClass} ${toast.isVisible ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0 pointer-events-none'}`}
      >
        <span>{toast.message}</span>
      </div>
    </div>
  );
}