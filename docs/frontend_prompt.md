# Loveable Prompt: MoneyLover App Replica

**Objective:** Build a fully functional financial management app similar to MoneyLover. The app should interact exclusively with the backend API I have attached. All UI elements, features, and workflows should reflect MoneyLover’s experience. 
This website should work in mobile as well as on web. IT should behave like a progressive web app also in mobile.




**Instructions:**

1. **Initialization:**  
   - On app start, initialize the frontend and ensure that it confirms the backend API is ready before rendering the UI. Using api/initdone
   - If initdone returns false, then give a page to pass value to api/init

2. **User Interface:**  
   - After initialization, display a UI similar to MoneyLover, including:  
     - Wallet overview ( very small on top with select option for wallet) 
     - Transaction lists  (this is main page, support sorting though transaction datetime, entry datetime and last modied datetime)
     - Categories management in a menu panel. Where we can add categories including their parents for non-root categories. 
     - Filters and summaries (daily, weekly, monthly, yearly)  and based on select datetime. Once selected show the total income and expenses on that selcted period
     - add tabs for switching transaction of previous or next period
     - Add/Edit/Delete transaction dialogs  
     - Add/Edit/Delete wallet dialogs  
   - Ensure smooth navigation between screens.

3. **Transactions:**  
   - Adding a transaction should allow selection of a category from a full **category tree**, date, amount, wallet, notes, and type (expense/income).  
   - Editing a transaction should prefill all details and allow updates.  
   - Deleting a transaction should prompt confirmation.
   -  Filters and summaries (daily, weekly, monthly, yearly) and also based on selection of start datetime and end datetime. Once selected show the total income and expenses on that selcted period
   - You need to date selector also to filter transaction, basically start date and end date
also for filtering the transactions
   - Keep UI minimal, Keep Daily, weekly, monthly, yealry and SelectDates in to single dropdown and once selected, same dropdown should show the selected time. This will save the space. 

4. **Categories:**  
   - Display the category tree when adding or editing transactions.  
   - Allow adding, updating, and deleting categories through API calls, along with adding parent category and wallets. 
   - Updating the category should have option to change the parent of categorym name, icon  etc.. You need to add opton to update parent for add category also 
   - Parent category need not be root category, it can be category of any level, because we support category in tree heirarchy 

5. **Wallets:**  
   - Display all wallets for a user.  
   - Allow creation, updating, and deletion of wallets using the API.  
   - Support assigning users to wallets and replacing wallets for a user.

6. **Users:**  
   - List, create, update, and delete users as provided by the backend API.  
   - Support assigning wallets to users and detaching them.

7. **API Integration:**  
   - All app actions (transactions, categories, wallets, users) must call the provided backend API.  
   - Use the API endpoints exactly as documented in the file I will provide.  
   - Handle API errors gracefully and show user-friendly messages.

8. **Additional Requirements:**  
   - The frontend should be dynamic and interactive.  
   - All CRUD operations should fully work.  
   - Include filtering, sorting, and searching transactions.  
   - Include proper UI feedback (loading indicators, success/error messages).  
   - Make sure the app feels and functions like MoneyLover.

**Note:** I attached a file containing all backend API endpoints and payload structures. Use this file as the single source of truth for API calls. Do not hardcode any data; all actions must be performed through the API.



Further Improvements: 
1. Wallet management and category management will go inside the Settings from the bottom
2. The sorting section shows "Transaction Date" along with the sort sign, the text can be removed from showing , texts should come only when someone goes and clicks in the sorting icon. 
3. Add caching of categories and other items 
4. Remember the last date and time of transaction so that same can be reused for next transaction ( in browser)
5. Sort based on entry date time if transaction date time is same 
6. Along with note of transaction, also show modified dataetime in very very tiny text , so that it is just visible
7. Adding caching for transactions, sync periodically and udpate the cache. Use local storage for cache. 
8. When click on a transaction then it open the edit windows . It should edit windows only when Long pressed on mobile and right click on chome. 
9. When clicked on edit transaction, it automatically opens the keyboard ( because first field is amount). 
10. While selecting category in transaction adding/update windows, the windows should go on top so that we can see the categories when putting text, otherwise categories are hidden by keyboard itself . Also make the category selection a little smaller, thin so that more categories can fit in the screen
11. The transaction add/edit window should go on top in mobile so that it is not covered by keyboard. Make things little compressed so that more can fit in. A big text boxes and all may not be needed but keep that UI better, it should not be too small also. 
12. Add sort by catogory along with transaction date, entry date and for each category show the transaction and also the total amount spent on that category can be shows beside the category name, 
13. In the transaction filer, add option for all transaction also along with daily, monthly, weekly, select dates. Keep the option name as "ALL" ( this would be specific to a wallet)
14. When searching for a cateogry, also include collasped subcategory if the filtered category is a parent category


Current Account: alexuidian@gmail.com

Further Improvement:

A. Add caching of categories, and refresh the cache when needed ( for eg. when user is updating the transaction, creating a new transaction or doing CRUD ops on Categories , do it in background)

B. Remember the time span selected for showing transaction 
   1. If is selected Monthly, weekly or daily then show current month, week, day
   2. If it is selected all, then show all
   3. if the custom dates are selected then keep selecting the custom dates

C. Even after caching, the data load is taking time.. just show the cached data and later update the live transactions if there is some change in the database

D. When searching for a cateogry while updating or creating a transaction, if category name I searched is parent category, then if I click on the expand button, It should show the sub categories also. Currently if my match is a parent category I do not see subcategories 

E. Like you added Modified time in each transaction, also add transaction datetime so that I can also see that. 

F. Category , ammount are still having padding, make it thinner so that more things can be accomodated

G. Note should be above Person/Payee when creating or updating a transaction

H. The save button on transaction is thinner than delete icon, make it equal 

I. Give option to select default wallet for a user

J. Give option for CRUD for Wallet Groups and when a wallet group is selected, you need to show transation for all wallets in that wallet group

K. When we are adding/updating a transaction , pop up windows apprears, on pressing back at that point, app itself it getting closed. It should only close the pop up and should not exit. 

L. Add another pop up if user is trying to exit the app using back button or some other button, asking if they really want to exit the app

M. If daaily is selected, the date should be shows as today/yesterday and then the actual date. It would be easier to ebter trnasaction . 

N. Remove the currency sign from everywhere. Keep it without $ sign. 