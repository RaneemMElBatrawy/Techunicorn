##Endpoint


1) POST "/register" - allow all users to register 

2) POST "/login" - allow users to login 

3) GET "/doctors" - view list of doctors. Access: Everyone can access. Admins can see everything. Doctors and Patients can only see firstname, lastname and email.

4) GET "/doctors/:id" - view doctor information. Access: Everyone can access. Admins can see everything. Doctors and Patients can only see firstname, lastname and email.

5) POST "/doctors/:id/slots" - view doctors available slots. Access: Everyone can access.

6) POST "/doctors/:id/book" - book an appointment with a doctor. Access: Only patients can access.

7) DELETE "/appointment/:id" - cancel an appointment. Access: admins and doctors who are booked for can delete any appointment.

8)POST "/doctors/availability/all" - view availability of all doctors. Access: Only admins and patients can acces. 

9) GET "/appointments/:id" - view any appointment details. Access: Admins can access but for doctors and patients who booked only, not all doctors and patients.

*) GET "/patients" - view list of patients. Access: Only admins can access this list.

*) GET "/patients/:id" - view lists of requested patient. Access: Only admins can access this list. 

10) GET "/patient/history" - view patient history. Access: Only the signed in patient but all admins and doctors can access. 

11) POST "/doctors/most/appointments" - view doctors with the most appointments in a given day . Access: Only admins.

12) POST "/doctors/most/hours" - view doctors who have 6+ hours total appointments in a day. Access; Only admins.
