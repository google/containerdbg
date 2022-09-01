FROM tomcat:8.5-jre11-openjdk

ARG PETCLINIC_WAR_APP=apps/petclinic.war

# Copy petclinic.war app
ADD ${PETCLINIC_WAR_APP} /usr/local/tomcat/webapps/petclinic.war

# Copy all relevant files from the directory tree of the Tomcat server.
ADD --chown=root:root tomcatServer.tar.gz /usr/local/tomcat
# ADD --chown=root:root file.txt /usr/local/tomcat

# Fixing CVE-2021-45960 in Tomcat community image
RUN apt update && apt upgrade -y libexpat1
