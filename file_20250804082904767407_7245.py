import json
import os

class AddressBook:
    def __init__(self, filename='address_book.json'):
        self.filename = filename
        self.contacts = self._load_contacts()
        self._next_id = self._get_next_id()

    def _get_next_id(self):
        if not self.contacts:
            return 1
        return max(contact['id'] for contact in self.contacts) + 1

    def _save_contacts(self):
        try:
            with open(self.filename, 'w') as f:
                json.dump(self.contacts, f, indent=4)
        except IOError as e:
            print(f"Error saving contacts: {e}")

    def _load_contacts(self):
        if not os.path.exists(self.filename):
            return []
        try:
            with open(self.filename, 'r') as f:
                return json.load(f)
        except json.JSONDecodeError:
            print("Error: Could not decode address book file. Starting with empty book.")
            return []
        except IOError as e:
            print(f"Error loading contacts: {e}")
            return []

    def add_contact(self, name, phone, email, address):
        contact = {
            'id': self._next_id,
            'name': name,
            'phone': phone,
            'email': email,
            'address': address
        }
        self.contacts.append(contact)
        self._next_id += 1
        self._save_contacts()
        return contact

    def view_contacts(self):
        return self.contacts

    def search_contacts(self, query):
        query = query.lower()
        results = [
            contact for contact in self.contacts
            if query in contact['name'].lower() or query in contact['phone'].lower()
        ]
        return results

    def delete_contact(self, contact_id):
        initial_len = len(self.contacts)
        self.contacts = [contact for contact in self.contacts if contact['id'] != contact_id]
        if len(self.contacts) < initial_len:
            self._save_contacts()
            return True
        return False

    def edit_contact(self, contact_id, new_name, new_phone, new_email, new_address):
        for contact in self.contacts:
            if contact['id'] == contact_id:
                contact['name'] = new_name
                contact['phone'] = new_phone
                contact['email'] = new_email
                contact['address'] = new_address
                self._save_contacts()
                return True
        return False

def display_contact(contact):
    print(f"ID: {contact['id']}")
    print(f"  Name: {contact['name']}")
    print(f"  Phone: {contact['phone']}")
    print(f"  Email: {contact['email']}")
    print(f"  Address: {contact['address']}")
    print("-" * 20)

def display_menu():
    print("\n--- Address Book CLI ---")
    print("1. Add Contact")
    print("2. View All Contacts")
    print("3. Search Contacts")
    print("4. Edit Contact")
    print("5. Delete Contact")
    print("6. Exit")
    print("------------------------")

def add_contact_cli(book):
    print("\n--- Add New Contact ---")
    name = input("Enter Name: ").strip()
    if not name:
        print("Name cannot be empty.")
        return
    phone = input("Enter Phone: ").strip()
    email = input("Enter Email: ").strip()
    address = input("Enter Address: ").strip()
    book.add_contact(name, phone, email, address)
    print("Contact added successfully!")

def view_contacts_cli(book):
    contacts = book.view_contacts()
    if not contacts:
        print("\nNo contacts in the address book.")
        return
    print("\n--- All Contacts ---")
    for contact in contacts:
        display_contact(contact)

def search_contacts_cli(book):
    query = input("\nEnter name or phone to search: ").strip()
    if not query:
        print("Search query cannot be empty.")
        return
    results = book.search_contacts(query)
    if not results:
        print(f"\nNo contacts found matching '{query}'.")
        return
    print(f"\n--- Search Results for '{query}' ---")
    for contact in results:
        display_contact(contact)

def edit_contact_cli(book):
    contact_id_str = input("\nEnter the ID of the contact to edit: ").strip()
    try:
        contact_id = int(contact_id_str)
    except ValueError:
        print("Invalid ID. Please enter a number.")
        return

    contact_to_edit = next((c for c in book.contacts if c['id'] == contact_id), None)
    if not contact_to_edit:
        print(f"No contact found with ID {contact_id}.")
        return

    print(f"\n--- Editing Contact ID: {contact_id} ({contact_to_edit['name']}) ---")
    print("Leave field blank to keep current value.")

    new_name = input(f"Enter new Name (current: {contact_to_edit['name']}): ").strip()
    new_phone = input(f"Enter new Phone (current: {contact_to_edit['phone']}): ").strip()
    new_email = input(f"Enter new Email (current: {contact_to_edit['email']}): ").strip()
    new_address = input(f"Enter new Address (current: {contact_to_edit['address']}): ").strip()

    name = new_name if new_name else contact_to_edit['name']
    phone = new_phone if new_phone else contact_to_edit['phone']
    email = new_email if new_email else contact_to_edit['email']
    address = new_address if new_address else contact_to_edit['address']

    if book.edit_contact(contact_id, name, phone, email, address):
        print("Contact updated successfully!")
    else:
        print("Failed to update contact (ID not found, though it should be).")

def delete_contact_cli(book):
    contact_id_str = input("\nEnter the ID of the contact to delete: ").strip()
    try:
        contact_id = int(contact_id_str)
    except ValueError:
        print("Invalid ID. Please enter a number.")
        return

    if book.delete_contact(contact_id):
        print(f"Contact with ID {contact_id} deleted successfully!")
    else:
        print(f"No contact found with ID {contact_id}.")

def main():
    book = AddressBook()

    while True:
        display_menu()
        choice = input("Enter your choice: ").strip()

        if choice == '1':
            add_contact_cli(book)
        elif choice == '2':
            view_contacts_cli(book)
        elif choice == '3':
            search_contacts_cli(book)
        elif choice == '4':
            edit_contact_cli(book)
        elif choice == '5':
            delete_contact_cli(book)
        elif choice == '6':
            print("Exiting Address Book. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-08-04 08:30:08
import json
import os

class Contact:
    def __init__(self, name, phone, email):
        self.name = name
        self.phone = phone
        self.email = email

    def __str__(self):
        return f"Name: {self.name}, Phone: {self.phone}, Email: {self.email}"

    def to_dict(self):
        return {
            "name": self.name,
            "phone": self.phone,
            "email": self.email
        }

    @classmethod
    def from_dict(cls, data):
        return cls(data["name"], data["phone"], data["email"])

class AddressBook:
    def __init__(self, filename="address_book.json"):
        self.filename = filename
        self.contacts = []
        self._load_contacts()

    def _load_contacts(self):
        if os.path.exists(self.filename):
            try:
                with open(self.filename, 'r') as f:
                    data = json.load(f)
                    self.contacts = [Contact.from_dict(item) for item in data]
            except json.JSONDecodeError:
                print("Error loading address book file. Starting with an empty book.")
                self.contacts = []
        else:
            self.contacts = []

    def _save_contacts(self):
        with open(self.filename, 'w') as f:
            json.dump([contact.to_dict() for contact in self.contacts], f, indent=4)

    def add_contact(self, name, phone, email):
        contact = Contact(name, phone, email)
        self.contacts.append(contact)
        self._save_contacts()
        print(f"Contact '{name}' added.")

    def list_contacts(self):
        if not self.contacts:
            print("Address book is empty.")
            return

        print("\n--- Your Contacts ---")
        for i, contact in enumerate(self.contacts):
            print(f"{i+1}. {contact}")
        print("---------------------\n")

    def search_contacts(self, term):
        term = term.lower()
        found_contacts = []
        for contact in self.contacts:
            if (term in contact.name.lower() or
                term in contact.phone.lower() or
                term in contact.email.lower()):
                found_contacts.append(contact)
        
        if not found_contacts:
            print(f"No contacts found matching '{term}'.")
            return []
        
        print(f"\n--- Search Results for '{term}' ---")
        for i, contact in enumerate(found_contacts):
            print(f"{i+1}. {contact}")
        print("----------------------------------\n")
        return found_contacts

    def edit_contact(self, index, new_name=None, new_phone=None, new_email=None):
        if 0 <= index < len(self.contacts):
            contact = self.contacts[index]
            if new_name is not None:
                contact.name = new_name
            if new_phone is not None:
                contact.phone = new_phone
            if new_email is not None:
                contact.email = new_email
            self._save_contacts()
            print(f"Contact '{contact.name}' updated.")
        else:
            print("Invalid contact number.")

    def delete_contact(self, index):
        if 0 <= index < len(self.contacts):
            deleted_contact = self.contacts.pop(index)
            self._save_contacts()
            print(f"Contact '{deleted_contact.name}' deleted.")
        else:
            print("Invalid contact number.")

def main():
    address_book = AddressBook()

    while True:
        print("\n--- Address Book Menu ---")
        print("1. Add Contact")
        print("2. List Contacts")
        print("3. Search Contacts")
        print("4. Edit Contact")
        print("5. Delete Contact")
        print("6. Exit")
        choice = input("Enter your choice: ")

        if choice == '1':
            name = input("Enter name: ")
            phone = input("Enter phone: ")
            email = input("Enter email: ")
            address_book.add_contact(name, phone, email)
        elif choice == '2':
            address_book.list_contacts()
        elif choice == '3':
            term = input("Enter search term: ")
            address_book.search_contacts(term)
        elif choice == '4':
            address_book.list_contacts()
            if not address_book.contacts:
                continue
            try:
                index = int(input("Enter the number of the contact to edit: ")) - 1
                if 0 <= index < len(address_book.contacts):
                    contact_to_edit = address_book.contacts[index]
                    print(f"Editing: {contact_to_edit}")
                    new_name = input(f"Enter new name (current: {contact_to_edit.name}, leave blank to keep): ")
                    new_phone = input(f"Enter new phone (current: {contact_to_edit.phone}, leave blank to keep): ")
                    new_email = input(f"Enter new email (current: {contact_to_edit.email}, leave blank to keep): ")
                    
                    address_book.edit_contact(
                        index,
                        new_name if new_name else contact_to_edit.name, # Pass current value if blank
                        new_phone if new_phone else contact_to_edit.phone,
                        new_email if new_email else contact_to_edit.email
                    )
                else:
                    print("Invalid contact number.")
            except ValueError:
                print("Invalid input. Please enter a number.")
        elif choice == '5':
            address_book.list_contacts()
            if not address_book.contacts:
                continue
            try:
                index = int(input("Enter the number of the contact to delete: ")) - 1
                address_book.delete_contact(index)
            except ValueError:
                print("Invalid input. Please enter a number.")
        elif choice == '6':
            print("Exiting Address Book. Goodbye!")
            break
        else:
            print("Invalid choice. Please try again.")

if __name__ == "__main__":
    main()