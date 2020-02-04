from astropy.constants import G, M_earth, R_earth
from astropy import units as u
import numpy as np

def cart_2_kep(r_vec,v_vec):
    #1
    h_bar = np.cross(r_vec,v_vec)
    h = np.linalg.norm(h_bar)
    #2
    r = np.linalg.norm(r_vec)
    v = np.linalg.norm(v_vec)
    #3
    E = 0.5*(v**2) - mu/r
    #4
    a = -mu/(2*E)
    #5
    e = np.sqrt(1 - (h**2)/(a*mu))
    #6
    i = np.arccos(h_bar[2]/h)
    #7
    omega_LAN = np.arctan2(h_bar[0],-h_bar[1])
    #8
    #beware of division by zero here
    lat = np.arctan2(np.divide(r_vec[2],(np.sin(i))),\
    (r_vec[0]*np.cos(omega_LAN) + r_vec[1]*np.sin(omega_LAN)))
    #9
    p = a*(1-e**2)
    nu = np.arctan2(np.sqrt(p/mu) * np.dot(r_vec,v_vec), p-r)
    #10
    omega_AP = lat - nu
    #11
    EA = 2*np.arctan(np.sqrt((1-e)/(1+e)) * np.tan(nu/2))
    #12
    n = np.sqrt(mu/(a**3))
    T = t - (1/n)*(EA - e*np.sin(EA))

    return a,e,i,omega_AP,omega_LAN,T, EA

def kep_2_cart(a,e,i,omega_AP,omega_LAN,T, EA):
    #1
    n = np.sqrt(mu/(a**3))
    M = n*(t - T)
    #2
    MA = EA - e*np.sin(EA)
    #3
    #
    # ERROR WAS HERE
   #nu = 2*np.arctan(np.sqrt((1-e)/(1+e)) * np.tan(EA/2))
    nu = 2*np.arctan(np.sqrt((1+e)/(1-e)) * np.tan(EA/2))
    #4
    r = a*(1 - e*np.cos(EA))
    #5
    h = np.sqrt(mu*a * (1 - e**2))
    #6
    Om = omega_LAN
    w =  omega_AP

    X = r*(np.cos(Om)*np.cos(w+nu) - np.sin(Om)*np.sin(w+nu)*np.cos(i))
    Y = r*(np.sin(Om)*np.cos(w+nu) + np.cos(Om)*np.sin(w+nu)*np.cos(i))
    Z = r*(np.sin(i)*np.sin(w+nu))

    #7
    p = a*(1-e**2)

    V_X = (X*h*e/(r*p))*np.sin(nu) - (h/r)*(np.cos(Om)*np.sin(w+nu) + \
    np.sin(Om)*np.cos(w+nu)*np.cos(i))
    V_Y = (Y*h*e/(r*p))*np.sin(nu) - (h/r)*(np.sin(Om)*np.sin(w+nu) - \
    np.cos(Om)*np.cos(w+nu)*np.cos(i))
    V_Z = (Z*h*e/(r*p))*np.sin(nu) - (h/r)*(np.cos(w+nu)*np.sin(i))

    return [X,Y,Z],[V_X,V_Y,V_Z]

mu = G.value*M_earth.value
Re = R_earth.value

#Test vectors
r_test = np.array([Re + 600.0*1000, 0, 50])
v_test = np.array([0, 6.5 * 1000, 0])
t = 0

a,e,i,omega_AP,omega_LAN,T, EA = cart_2_kep(r_test,v_test)
r_test2, v_test2 = kep_2_cart(a,e,i,omega_AP,omega_LAN,T, EA)

print('foo: {0}\n'.format(Re))
print('a: {0}\ne: {1}\ni: {2}\nΩ: {3}\nω: {4}\nf: {5}\n'.format(a,e,i,omega_AP,omega_LAN,T))

print(r_test2)
print(v_test2)
#print(r_test2 - r_test)
#print(v_test2 - v_test)
